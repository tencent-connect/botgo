// Package client 默认的 websocket client 实现。
package client

import (
	"encoding/json"
	"fmt"
	"time"

	wss "github.com/gorilla/websocket" // 是一个流行的 websocket 客户端，服务端实现
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/websocket"
)

// DefaultQueueSize 监听队列的缓冲长度
const DefaultQueueSize = 10000

// Setup 依赖注册
func Setup() {
	websocket.Register(&Client{})
}

// New 新建一个连接对象
func (c *Client) New(session dto.Session) websocket.WebSocket {
	return &Client{
		messageQueue:    make(messageChan, DefaultQueueSize),
		session:         &session,
		closeChan:       make(closeErrorChan, 10),
		heartBeatTicker: time.NewTicker(60 * time.Second), // 先给一个默认 ticker，在收到 hello 包之后，会 reset
	}
}

// Client websocket 连接客户端
type Client struct {
	version         int
	conn            *wss.Conn
	messageQueue    messageChan
	session         *dto.Session
	user            *dto.WSUser
	closeChan       closeErrorChan
	heartBeatTicker *time.Ticker // 用于维持定时心跳
}

type messageChan chan []byte
type closeErrorChan chan error

// Connect 连接到 websocket
func (c *Client) Connect() error {
	if c.session.URL == "" {
		return errs.ErrURLInvalid
	}

	var err error
	c.conn, _, err = wss.DefaultDialer.Dial(c.session.URL, nil)
	if err != nil {
		log.Errorf("%s, connect err: %v", c.session, err)
		return err
	}
	log.Infof("%s, url %s, connected", c.session, c.session.URL)

	return nil
}

// Listening 开始监听，会阻塞进程，内部会从事件队列不断的读取事件，解析后投递到注册的 event handler，如果读取消息过程中发生错误，会循环
// 定时心跳也在这里维护
func (c *Client) Listening() error {
	defer c.Close()
	// reading message
	go c.readMessageToQueue()
	// read message from queue and handle,in goroutine to avoid business logic block closeChan and heartBeatTicker
	go c.listenMessageAndHandle()

	// handler message
	for {
		select {
		case err := <-c.closeChan:
			// 关闭连接的错误码 https://bot.q.qq.com/wiki/develop/api/gateway/error/error.html
			log.Errorf("%s Listening stop. err is %v", c.session, err)
			// 不能够 identify 的错误
			if wss.IsCloseError(err, 4914, 4915) {
				return errs.New(errs.CodeConnCloseCantIdentify, err.Error())
			}
			// 这里用 UnexpectedCloseError，如果有需要排除在外的 close error code，可以补充在第二个参数上
			// 4009: session time out, 发了 reconnect 之后马上关闭连接时候的错误码，这个是允许 resume 的
			if wss.IsUnexpectedCloseError(err, 4009) {
				return errs.New(errs.CodeConnCloseCantResume, err.Error())
			}
			return err
		case <-c.heartBeatTicker.C:
			log.Debugf("%s listened heartBeat", c.session)
			heartBeatEvent := &dto.WSPayload{
				WSPayloadBase: dto.WSPayloadBase{
					OPCode: dto.WSHeartbeat,
				},
				Data: c.session.LastSeq,
			}
			// 不处理错误，Write 内部会处理，如果发生发包异常，会通知主协程退出
			_ = c.Write(heartBeatEvent)
		}
	}
}

func (c *Client) listenMessageAndHandle() {
	defer func() {
		// panic，一般是由于业务自己实现的 handle 不完善导致
		// 打印日志后，关闭这个连接，进入重连流程
		if err := recover(); err != nil {
			log.Errorf("%s panic err: %v", c.session, err)
			c.closeChan <- fmt.Errorf("panic: %v", err)
		}
	}()
	for message := range c.messageQueue {
		log.Debugf("%s listened message", c.session)
		event := &dto.WSPayload{}
		if err := json.Unmarshal(message, event); err != nil {
			log.Errorf("%s json failed, %v", c.session, err)
			continue
		}
		c.writeSeq(event.Seq)
		// 处理内置的一些事件，如果处理成功，则这个事件不再投递给业务
		if c.isHandleBuildIn(event, message) {
			continue
		}
		// ready 事件需要特殊处理
		if event.Type == "READY" {
			c.readyHandler(message)
			continue
		}
		// 解析具体事件，并投递给业务注册的 handler
		if err := parseAndHandle(event, message); err != nil {
			log.Errorf("%s parseAndHandle failed, %v", c.session, err)
		}
	}
	log.Infof("%s message queue is closed", c.session)
}

func (c *Client) Write(message *dto.WSPayload) error {
	m, _ := json.Marshal(message)
	log.Infof("%s write message, %v", c.session, string(m))

	if err := c.conn.WriteMessage(wss.TextMessage, m); err != nil {
		log.Errorf("%s WriteMessage failed, %v", c.session, err)
		c.closeChan <- err
		return err
	}
	return nil
}

// Resume 重连
func (c *Client) Resume() error {
	event := &dto.WSPayload{
		Data: &dto.WSResumeData{
			Token:     c.session.Token.GetString(),
			SessionID: c.session.ID,
			Seq:       c.session.LastSeq,
		},
	}
	event.OPCode = dto.WSResume // 内嵌结构体字段，单独赋值
	return c.Write(event)
}

// Identify 对一个连接进行鉴权，并声明监听的 shard 信息
func (c *Client) Identify() error {
	// 避免传错 intent
	if c.session.Intent == 0 {
		c.session.Intent = dto.IntentGuilds
	}
	event := &dto.WSPayload{
		Data: &dto.WSIdentityData{
			Token:   c.session.Token.GetString(),
			Intents: c.session.Intent,
			Shard: []uint32{
				c.session.Shards.ShardID,
				c.session.Shards.ShardCount,
			},
		},
	}
	event.OPCode = dto.WSIdentity
	return c.Write(event)
}

// Close 关闭连接
func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		log.Errorf("%s, close conn err: %v", c.session, err)
	}
	c.heartBeatTicker.Stop()
}

// Session 获取client的session信息
func (c *Client) Session() *dto.Session {
	return c.session
}

// isHandleBuildIn 内置的事件处理，处理那些不需要业务方处理的事件
// return true 的时候说明事件已经被处理了
func (c *Client) isHandleBuildIn(event *dto.WSPayload, message []byte) bool {
	switch event.OPCode {
	case dto.WSHello: // 接收到 hello 后需要开始发心跳
		c.startHeartBeatTicker(message)
		return true
	case dto.WSHeartbeatAck: // 心跳 ack 不需要业务处理
		return true
	case dto.WSReconnect: // 达到连接时长，需要重新连接，此时可以通过 resume 续传原连接上的事件
		c.closeChan <- errs.ErrNeedReConnect
		return true
	case dto.WSInvalidSession: // 无效的 sessionLog，需要重新鉴权
		c.closeChan <- errs.ErrInvalidSession
		return true
	default:
		return false
	}
}

// startHeartBeatTicker 启动定时心跳
func (c *Client) startHeartBeatTicker(message []byte) {
	helloData := &dto.WSHelloData{}
	if err := parseData(message, helloData); err != nil {
		log.Errorf("%s hello data parse failed, %v, message %v", c.session, err, message)
	}
	// 根据 hello 的回包，重新设置心跳的定时器时间
	c.heartBeatTicker.Reset(time.Duration(helloData.HeartbeatInterval) * time.Millisecond)
}

func (c *Client) readMessageToQueue() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Errorf("%s read message failed, %v, message %s", c.session, err, string(message))
			close(c.messageQueue)
			c.closeChan <- err
			return
		}
		log.Infof("%s receive message, %s", c.session, string(message))
		c.messageQueue <- message
	}
}

// readyHandler 针对ready返回的处理，需要记录 sessionID 等相关信息
func (c *Client) readyHandler(message []byte) {
	readyData := &dto.WSReadyData{}
	if err := parseData(message, readyData); err != nil {
		log.Errorf("%s parseReadyData failed, %v, message %v", c.session, err, message)
	}
	c.version = readyData.Version
	// 基于 ready 事件，更新 session 信息
	c.session.ID = readyData.SessionID
	c.session.Shards.ShardID = readyData.Shard[0]
	c.session.Shards.ShardCount = readyData.Shard[1]
	c.user = &dto.WSUser{
		ID:       readyData.User.ID,
		Username: readyData.User.Username,
		Bot:      readyData.User.Bot,
	}
}

func (c *Client) writeSeq(seq uint32) {
	if seq > 0 {
		c.session.LastSeq = seq
	}
}
