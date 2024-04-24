// Package remote 基于 redis list 实现的分布式 session manager。
package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/sessions/manager"
	"github.com/tencent-connect/botgo/sessions/remote/lock"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

const (
	// 分布式锁的默认key，可以从外部通过 option 来指定
	defaultClusterKey = "defaultCluster"
	// session 队列的key后缀，实际上的key为 `fmt.Sprintf("%s_%s", r.clusterKey, sessionQueueSuffix)`
	sessionQueueSuffix = "sessionsQueue"
	// 分发shard的实例的分布式锁的默认过期时间
	distributeLockExpireTime = 60 * time.Second
	// 每个不同的shard实例的分布式锁，用于避免同个 shard 被启动多个实例
	shardLockExpireTime = 30 * time.Second
)

// RedisManager 基于 redis 的 session 管理器，实现分布式 websocket 监听
type RedisManager struct {
	clusterKey         string
	sessionQueueKey    string
	client             *redis.Client
	sessionProduceChan chan dto.Session // 抢到锁的服务，用于持续生产session到redis list的本地chan
}

// New 创建一个新的基于 redis 的 session 管理器
// 使用 go-redis 调用 redis，超时时间请在 NewClient 时候设置
func New(client *redis.Client, opts ...Option) *RedisManager {
	r := &RedisManager{
		clusterKey: defaultClusterKey,
		client:     client,
	}
	for _, opt := range opts {
		opt(r)
	}
	// 针对不同的分布式key，设置不同的 queue key
	r.sessionQueueKey = fmt.Sprintf("%s_%s", r.clusterKey, sessionQueueSuffix)
	return r
}

// Start 启动 redis 的 session 管理器
func (r *RedisManager) Start(apInfo *dto.WebsocketAP, token *token.Token, intents *dto.Intent) error {
	defer log.Sync()
	if err := manager.CheckSessionLimit(apInfo); err != nil {
		log.Errorf("[ws/session/redis] session limited apInfo: %+v", apInfo)
		return err
	}
	startInterval := manager.CalcInterval(apInfo.SessionStartLimit.MaxConcurrency)
	log.Infof("[ws/session/redis] will start %d sessions and per session start interval is %s",
		apInfo.Shards, startInterval)

	// session 生产队列
	r.sessionProduceChan = make(chan dto.Session, apInfo.Shards)

	// 进行初始的session分发，抢锁，分发
	// 锁60s，抢到锁的进程，需要每30s续期一次，只要自己还存活，就不能够让另外的进程抢到锁重新进行shards分发
	ctx := context.Background()
	distributeLock := lock.New(r.clusterKey, uuid.New().String(), r.client)
	if err := distributeLock.Lock(ctx, distributeLockExpireTime); err == nil {
		log.Infof("[ws/session/redis] got distribute lock! i will do distributeSession, key: %s", r.clusterKey)
		// 抢到锁的进行初次分发
		if err = r.distributeSession(apInfo, token, intents); err != nil {
			log.Errorf("[ws/session/redis] distribute sessions failed: %v", err)
			return err
		}
		go distributeLock.StartRenew(ctx, distributeLockExpireTime)
	} else {
		log.Errorf("got lock failed, err: %v", err)
	}

	// 持续 produce session，遇到网络问题在 chan 中重试
	// 对于抢到了锁的服务，生产第一批session到redis list
	// 对于没有抢到锁的服务，当ws异常，把session放回到 redis list 中，重新分发
	go r.sessionProducer(startInterval)

	return r.consume(startInterval)
}

func (r *RedisManager) consume(startInterval time.Duration) error {
	log.Debug("[ws/session/redis] start consume for session")
	for {
		// brpop 返回 key value
		data, err := r.client.BRPop(context.Background(), startInterval*2, r.sessionQueueKey).Result()
		if err != nil {
			if err != redis.Nil {
				log.Errorf("[ws/session/redis] rpop failed, err: %v", err)
			}
			continue
		}
		if len(data) < 2 {
			log.Errorf("[ws/session/redis] data is not valid, data: %+v", data)
			continue
		}
		log.Debugf("[ws/session/redis] consume data: %s", data)

		session := &dto.Session{}
		if err := json.Unmarshal([]byte(data[1]), session); err != nil {
			// 解析出错，不放回去，直接丢弃
			log.Errorf("[ws/session/redis] unmarshal session failed, err: %v", err)
			continue
		}

		go r.newConnect(*session)
		time.Sleep(startInterval) // 启动一个连接后，等待一下，避免触发服务端的并发控制
	}
}

// getShardLockKey 获取 shard 的锁
func (r *RedisManager) getShardLockKey(session dto.Session) string {
	return fmt.Sprintf("%s_shard_%d_%d",
		r.clusterKey, session.Shards.ShardID, session.Shards.ShardCount)
}

// newConnect 启动一个新的连接，如果连接在监听过程中报错了，或者被远端关闭了链接，需要识别关闭的原因，能否继续 resume
// 如果能够 resume，则往 sessionChan 中放入带有 sessionID 的 session
// 如果不能，则清理掉 sessionID，将 session 放入 sessionChan 中
// session 的启动，交给 start 中的 for 循环执行，session 不自己递归进行重连，避免递归深度过深
func (r *RedisManager) newConnect(session dto.Session) {
	ctx := context.Background()
	// 锁 shard，避免针对相同 shard 消费重复了
	shardLock := lock.New(r.getShardLockKey(session), uuid.NewString(), r.client)
	if err := shardLock.Lock(ctx, shardLockExpireTime); err != nil {
		// shard 抢锁失败，把 session 放回去，避免上一个 session 的锁释放失败，导致下一个 session 无法启动
		r.sessionProduceChan <- session
		return
	}
	go shardLock.StartRenew(ctx, shardLockExpireTime)

	wsClient := websocket.ClientImpl.New(session)
	if err := wsClient.Connect(); err != nil {
		log.Error(err)
		r.sessionProduceChan <- session // 连接失败，丢回去队列排队重连
		return
	}
	var err error
	// 如果 session id 不为空，则执行的是 resume 操作，如果为空，则执行的是 identify 操作
	if session.ID != "" {
		err = wsClient.Resume()
	} else {
		// 初次鉴权
		err = wsClient.Identify()
	}
	if err != nil {
		log.Errorf("[ws/session/remote] Identify/Resume err %+v", err)
		return
	}
	if err := wsClient.Listening(); err != nil {
		log.Errorf("[ws/session/remote] Listening err %+v", err)
		currentSession := wsClient.Session()
		// 对于不能够进行重连的session，需要清空 session id 与 seq
		if manager.CanNotResume(err) {
			currentSession.ID = ""
			currentSession.LastSeq = 0
		}
		// 一些错误不能够鉴权，比如机器人被封禁，这里就直接退出了
		if manager.CanNotIdentify(err) {
			msg := fmt.Sprintf("can not identify because server return %+v, so process exit", err)
			log.Errorf(msg)
			panic(msg) // 当机器人被下架，或者封禁，将不能再连接，所以 panic
		}
		// 将 session 放到 session chan 中，用于启动新的连接，释放锁，当前连接退出
		shardLock.StopRenew()
		if err := shardLock.Release(ctx); err != nil {
			log.Errorf("[ws/session/remote] release shardLock failed, err: %s", err)
		}
		r.sessionProduceChan <- *currentSession
		return
	}
}
