# 分布式 session manager

这是一个基于 `redis` 的 `list` 数据结构的分布式 session manager。

## 实现原理

1.基于 redis 实现的分布式锁，启动的时候先抢锁，抢到锁的服务实例根据从 openapi 拉取到的 shards 进行 session 的分发

2.启动一个本地的 `sessionProduceChan` 的消费，用于将需重连的 session 重新 push 到 redis 中，如果 push 失败，放回 chan 进行下一次重试

3.启动一个消费者，从 redis pop 数据，解析 session，然后创建新的 websocket client 连接

4.如果在处理 websocket 数据过程中出现连接错误等情况，将 session 放回到 `sessionProduceChan` 中，重新进行分发 

## 并发控制

由于服务端对于同时连接的 websocket 连接有并发限制，所以从 `sessionProduceChan` 拿到一个 session push 到 redis 之前，会等待一个并发间隔

在创建了一个新的 websocket 连接时候，也会等待一个时间间隔

## 使用方法

[参考代码](../../testcase/redis_session_manager_test.go)
