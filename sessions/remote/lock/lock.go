// Package lock 一个基于 redis 的分布式锁实现。
package lock

import (
	"context"
	"errors"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/tencent-connect/botgo/log"
)

// ErrorNotOk redis 写失败
var ErrorNotOk = errors.New("redis write not ok")

// Lock 一个基于redis的锁实现
type Lock struct {
	lockKey       string
	lockValue     string
	client        *redis.Client
	renewTicker   *time.Ticker // 用于续期的ticker，默认为超时时间的 1/3
	stopRenewChan chan bool    // 用于停止 renew
}

// New 创建一个锁
func New(key, value string, client *redis.Client) *Lock {
	return &Lock{
		lockKey:   key,
		lockValue: value,
		client:    client,
	}
}

// Lock 加锁
func (l *Lock) Lock(ctx context.Context, expire time.Duration) error {
	success, err := l.client.SetNX(ctx, l.lockKey, l.lockValue, expire).Result()
	if err != nil {
		return err
	}
	if !success {
		return ErrorNotOk
	}
	return nil
}

// StartRenew 开始续期任务，需要放到 goroutine 中执行
func (l *Lock) StartRenew(ctx context.Context, expire time.Duration) {
	if expire == 0 {
		return
	}
	l.stopRenewChan = make(chan bool)
	l.renewTicker = time.NewTicker(expire / 3)
	defer l.renewTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Infof("[lock] context done, stop renew, %+v", l)
			return
		case <-l.stopRenewChan:
			log.Infof("[lock] renew stop, %+v", l)
			return
		case <-l.renewTicker.C:
			if err := l.Renew(ctx, expire); err != nil {
				log.Errorf("[lock] renew lock failed, lock: %+v, err: %v", l, err)
				continue
			}
			log.Debugf("[lock] renew lock ok, lock: %+v", l)
		}
	}
}

// StopRenew 停掉续期
func (l *Lock) StopRenew() {
	if l.stopRenewChan == nil {
		return
	}
	l.stopRenewChan <- true
}

// Renew 续期锁
func (l *Lock) Renew(ctx context.Context, expire time.Duration) error {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("expire", KEYS[1], ARGV[2])
	else
		return 0
	end
	`
	renew := redis.NewScript(script)
	return renew.Run(ctx, l.client, []string{l.lockKey}, l.lockValue, expire.Seconds()).Err()
}

// Release 释放锁
func (l *Lock) Release(ctx context.Context) error {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	release := redis.NewScript(script)
	return release.Run(ctx, l.client, []string{l.lockKey}, l.lockValue).Err()
}
