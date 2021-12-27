package lock

import (
	"context"
	"testing"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func TestLock_Lock(t *testing.T) {
	conn := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DialTimeout:  800 * time.Millisecond,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	key := "lock-key-test"
	value := "ddd"
	expire := 4 * time.Second
	ctx := context.Background()

	lock := New(key, value, conn)

	t.Run("lock", func(t *testing.T) {
		if err := lock.Lock(ctx, expire); err != nil {
			t.Error(err)
		}
	})
	t.Run("renew", func(t *testing.T) {
		if err := lock.Renew(ctx, expire); err != nil {
			t.Error(err)
		}
	})
	t.Run("release", func(t *testing.T) {
		if err := lock.Release(ctx); err != nil {
			t.Error(err)
		}
	})

	t.Run("renew goroutine and check", func(t *testing.T) {
		if err := lock.Lock(ctx, expire); err != nil {
			t.Error(err)
		}
		go lock.StartRenew(ctx, expire)
		time.Sleep(expire + 2*time.Second)
		// renew 持续在跑，这里不应该再抢到锁
		if err := lock.Lock(ctx, expire); err == nil {
			t.Error("want lock err, but got nil")
		}
		lock.StopRenew()
		time.Sleep(expire + 2*time.Second)
		if err := lock.Lock(ctx, expire); err != nil {
			t.Error(err)
		}
	})
}
