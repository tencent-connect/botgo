// Package local 基于 golang chan 实现的单机 manager。
package local

import (
	"fmt"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/sessions/manager"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

// New 创建本地session管理器
func New() *ChanManager {
	return &ChanManager{}
}

// ChanManager 默认的本地 session manager 实现
type ChanManager struct {
	sessionChan chan dto.Session
}

// Start 启动本地 session manager
func (l *ChanManager) Start(apInfo *dto.WebsocketAP, token *token.Token, intents *dto.Intent) error {
	defer log.Sync()
	if err := manager.CheckSessionLimit(apInfo); err != nil {
		log.Errorf("[ws/session/local] session limited apInfo: %+v", apInfo)
		return err
	}
	startInterval := manager.CalcInterval(apInfo.SessionStartLimit.MaxConcurrency)
	log.Infof("[ws/session/local] will start %d sessions and per session start interval is %s",
		apInfo.Shards, startInterval)

	// 按照shards数量初始化，用于启动连接的管理
	l.sessionChan = make(chan dto.Session, apInfo.Shards)
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
		l.sessionChan <- session
	}

	for session := range l.sessionChan {
		// MaxConcurrency 代表的是每 5s 可以连多少个请求
		time.Sleep(startInterval)
		go l.newConnect(session)
	}
	return nil
}

// newConnect 启动一个新的连接，如果连接在监听过程中报错了，或者被远端关闭了链接，需要识别关闭的原因，能否继续 resume
// 如果能够 resume，则往 sessionChan 中放入带有 sessionID 的 session
// 如果不能，则清理掉 sessionID，将 session 放入 sessionChan 中
// session 的启动，交给 start 中的 for 循环执行，session 不自己递归进行重连，避免递归深度过深
func (l *ChanManager) newConnect(session dto.Session) {
	defer func() {
		// panic 留下日志，放回 session
		if err := recover(); err != nil {
			websocket.PanicHandler(err, &session)
			l.sessionChan <- session
		}
	}()
	wsClient := websocket.ClientImpl.New(session)
	if err := wsClient.Connect(); err != nil {
		log.Error(err)
		l.sessionChan <- session // 连接失败，丢回去队列排队重连
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
		log.Errorf("[ws/session] Identify/Resume err %+v", err)
		return
	}
	if err := wsClient.Listening(); err != nil {
		log.Errorf("[ws/session] Listening err %+v", err)
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
		// 将 session 放到 session chan 中，用于启动新的连接，当前连接退出
		l.sessionChan <- *currentSession
		return
	}
}
