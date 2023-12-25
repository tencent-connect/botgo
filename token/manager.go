// Package token 用于调用 openapi，websocket 的 token 对象。
package token

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"encoding/json"

	"github.com/tencent-connect/botgo/log"
	"gopkg.in/yaml.v3"
)

// Type token 类型
type Type string

// TokenType
const (
	TypeNormal Type = "Bearer"
	TypeQQBot  Type = "QQBot"
)

type ManagerState int

const (
	ManagerStateUninitialized ManagerState = 0
	ManagerStateWorking       ManagerState = 1
	ManagerStateStopped       ManagerState = 2
)

// ITokenManager token manager 接口定义
type ITokenManager interface {
	Init(ctx context.Context) (err error)
	Close()
	State() ManagerState
	GetAccessToken() *AccessToken
	doRefreshToken() error
	GetRefreshSigCh() chan interface{}
}

// Manager token manager
type Manager struct {
	appID             uint64
	appSecret         string
	Type              Type
	token             *AccessToken
	lock              sync.RWMutex
	forceRefreshToken chan interface{}
	closeCh           chan int
	once              sync.Once
}

type tokenData struct {
	AppID        uint64
	ClientSecret string
	Type         Type
	Token        *AccessToken
}

// GetRefreshSigCh 获取同步刷新信号的 channel
func (m *Manager) GetRefreshSigCh() chan interface{} {
	return m.forceRefreshToken
}

// UnmarshalJSON 反序列化 json
func (m *Manager) UnmarshalJSON(data []byte) error {
	conf := &tokenData{}
	err := json.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	m.appID = conf.AppID
	m.appSecret = conf.ClientSecret
	m.Type = conf.Type
	m.token = conf.Token
	return nil
}

// MarshalJSON 序列化 json
func (m *Manager) MarshalJSON() ([]byte, error) {
	return json.Marshal(tokenData{
		AppID:        m.appID,
		ClientSecret: m.appSecret,
		Type:         m.Type,
		Token:        m.token,
	})
}

// NewManager 创建一个新的 Manager
func NewManager(tokenType Type) *Manager {
	return &Manager{
		Type: tokenType,
	}
}

// NewBotTokenManager 机器人身份的 token
func NewBotTokenManager(appID uint64, secret string) *Manager {
	manager := &Manager{
		appID:     appID,
		appSecret: secret,
		Type:      TypeQQBot,
	}
	return manager
}

// LoadAppAccFromYAML 从配置中读取 appid 和 token
func (m *Manager) LoadAppAccFromYAML(file string) error {
	var conf struct {
		AppID uint64 `yaml:"appid"`
		Token string `yaml:"token"`
	}
	content, err := os.ReadFile(file)
	if err != nil {
		log.Errorf("read token from file failed, err: %v", err)
		return err
	}
	if err = yaml.Unmarshal(content, &conf); err != nil {
		log.Errorf("parse config failed, err: %v", err)
		return err
	}
	m.appID = conf.AppID
	m.appSecret = conf.Token
	return nil
}

// Init 初始化，开始定时刷新token
func (m *Manager) Init(ctx context.Context) (err error) {
	m.once.Do(func() {
		m.forceRefreshToken = make(chan interface{}, 10)
		m.closeCh = make(chan int, 1)
		err = startRefreshAccessToken(ctx, m)
	})
	return err
}

// Close 停止定时刷新token
func (m *Manager) Close() {
	if isChanClose(m.closeCh) {
		return
	}
	close(m.closeCh)
}

func isChanClose(ch chan int) bool {
	select {
	case _, received := <-ch:
		return !received
	default:
	}
	return false
}

// GetAppID 取得Token中的appid
func (m *Manager) GetAppID() uint64 {
	return m.appID
}

// GetTokenValue 获取授权头字符串
func (m *Manager) GetTokenValue() string {
	return fmt.Sprintf("%s %s", m.Type, m.GetAccessToken().GetToken())
}

// GetAccessToken 取得鉴权Token
func (m *Manager) GetAccessToken() *AccessToken {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.token
}
func (m *Manager) setAccessToken(accessToken *AccessToken) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.token = accessToken
}

func (m *Manager) doRefreshToken() error {
	tokenInfo, err := retrieveToken(fmt.Sprint(m.appID), m.appSecret)
	if err != nil {
		return err
	}
	m.setAccessToken(tokenInfo)
	return nil
}

// State 判断Manager是否停止
func (m *Manager) State() ManagerState {
	if m.closeCh == nil {
		return ManagerStateUninitialized
	}
	if isChanClose(m.closeCh) {
		return ManagerStateStopped
	}
	return ManagerStateWorking
}

// startRefreshAccessToken 启动获取AccessToken的后台刷新
func startRefreshAccessToken(ctx context.Context, m ITokenManager) (err error) {
	if err = m.doRefreshToken(); err != nil {
		return
	}
	tokenExpireTime := m.GetAccessToken().UpdateTime.Add(time.Duration(m.GetAccessToken().ExpiresIn) * time.Second)
	tokenTTL := getTokenTTL(tokenExpireTime.Sub(time.Now()).Seconds())
	go func() {
		for {
			//如果manager已经停止，协程退出
			if m.State() != ManagerStateWorking {
				return
			}
			select {
			case <-time.NewTimer(time.Duration(tokenTTL) * time.Second).C:
			case reason := <-m.GetRefreshSigCh():
				log.Infof("force refresh token, reason:%v", reason)
			case <-ctx.Done():
				log.Warnf("recv ctx:%v exit refresh token", ctx.Err())
				m.Close()
				return
			}
			if err = m.doRefreshToken(); err != nil {
				log.Errorf("refresh access token failed:%v", err)
				tokenTTL = minTimeGap
				continue
			}
			tokenExpireTime = m.GetAccessToken().UpdateTime.Add(
				time.Duration(m.GetAccessToken().ExpiresIn) * time.Second)
			tokenTTL = getTokenTTL(tokenExpireTime.Sub(time.Now()).Seconds())
			log.Infof("tokenTTL:%d", tokenTTL)
		}
	}()
	return err
}
