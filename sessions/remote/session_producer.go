package remote

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/token"
)

// distributeSession 根据 shards 生产初始化的 session，这里需要抢一个分布式锁，抢到锁的服务器，负责把session都生产到 redis 中
func (r *RedisManager) distributeSession(apInfo *dto.WebsocketAP, token *token.Token, intents *dto.Intent) error {
	// clear，报错也不影响
	if err := r.client.Del(context.Background(), r.sessionQueueKey); err != nil {
		log.Errorf("[ws/session/redis] clear session list failed: %v", err)
	}
	for i := uint32(0); i < apInfo.Shards; i++ {
		session := dto.Session{
			URL:     apInfo.URL,
			Token:   *token,
			Intent:  *intents,
			LastSeq: 0,
			Shards: dto.ShardConfig{
				ShardID:    i,
				ShardCount: apInfo.Shards,
			},
		}
		r.sessionProduceChan <- session
	}

	return nil
}

// sessionProducer 从 chan 取到session，push 到 redis，push 失败放回 chan
func (r *RedisManager) sessionProducer(startInterval time.Duration) {
	for session := range r.sessionProduceChan {
		time.Sleep(startInterval) // 每次生产需要等待一个间隔，控制消费者连接并发
		if err := r.produce(session); err != nil {
			log.Errorf("[ws/session/redis] produce session failed: %v", err)
			r.sessionProduceChan <- session // 放回去重试
		}
	}
}

func (r *RedisManager) produce(session dto.Session) error {
	data, err := json.Marshal(session)
	log.Debugf("[ws][session/redis] produce session data is %s", string(data))
	if err != nil {
		return ErrSessionMarshalFailed
	}
	return r.client.LPush(context.Background(), r.sessionQueueKey, data).Err()
}
