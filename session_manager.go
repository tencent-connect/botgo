package botgo

import (
	"math"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

// defaultSessionManager 默认实现的 session manager 为单机版本
// 如果业务要自行实现分布式的 session 管理，则实现 SessionManger 后替换掉 defaultSessionManager
var defaultSessionManager SessionManager = &localSession{}

// SessionManager 接口，管理session
type SessionManager interface {
	// Start 启动连接
	Start(apInfo *dto.WebsocketAP, token *token.Token, intents *dto.Intent) error
}

// CanNotResumeErrSet 不能进行 resume 操作的错误码
var CanNotResumeErrSet = map[int]bool{
	errs.CodeConnCloseErr:   true,
	errs.CodeInvalidSession: true,
}

// CanNotResume 是否是不能够 resume 的错误
func CanNotResume(err error) bool {
	e := errs.Error(err)
	if flag, ok := CanNotResumeErrSet[e.Code()]; ok {
		return flag
	}
	return false
}

// CalcInterval 根据并发要求，计算连接启动间隔
func CalcInterval(maxConcurrency uint32) time.Duration {
	f := math.Round(5 / float64(maxConcurrency))
	return time.Duration(f) * time.Second
}

// localSession 默认的本地 session manager 实现
type localSession struct {
	sessionChan chan dto.Session
}

func (l *localSession) Start(apInfo *dto.WebsocketAP, token *token.Token, intents *dto.Intent) error {
	if err := l.checkSessionLimit(apInfo); err != nil {
		log.Errorf("[ws/session]session limited apInfo: %+v", apInfo)
		return err
	}
	startInterval := CalcInterval(apInfo.SessionStartLimit.MaxConcurrency)
	log.Infof("[ws/session] will start %d sessions and per session start interval is %s",
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
func (l *localSession) newConnect(session dto.Session) {
	wsClient := websocket.ClientImpl.New(session)
	if err := wsClient.Connect(); err != nil {
		log.Error(err)
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
		if CanNotResume(err) {
			currentSession.ID = ""
			currentSession.LastSeq = 0
		}
		// 将 session 放到 session chan 中，用于启动新的连接，当前连接退出
		l.sessionChan <- *currentSession
		return
	}
}

// checkSessionLimit 检查链接数是否达到限制
func (l *localSession) checkSessionLimit(apInfo *dto.WebsocketAP) error {
	if apInfo.Shards > apInfo.SessionStartLimit.Remaining {
		return errs.ErrSessionLimit
	}
	return nil
}
