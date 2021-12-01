package remote

// Option is a function that configures a Remote.
type Option func(manager *RedisManager)

// WithClusterKey 自定义集群key，用于创建分布式锁与redis list
func WithClusterKey(key string) Option {
	return func(m *RedisManager) {
		m.clusterKey = key
	}
}
