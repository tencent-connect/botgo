// Package token 用于调用 openapi，websocket 的 token 对象。
package token

import (
	"context"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/tencent-connect/botgo/log"
)

// Type token 类型
type Type string

// TokenType
const (
	TypeBot    Type = "Bot"
	TypeNormal Type = "Bearer"
	TypeQQBot  Type = "QQBot"
)

// Token 用于调用接口的 token 结构
type Token struct {
	appID        uint64
	clientSecret string
	tokenURL     string
	Type         Type
	authToken    *AuthTokenInfo
}

// SetTokenURL 设置获取accessToken的调用地址
func (t *Token) SetTokenURL(tokenURL string) {
	t.tokenURL = tokenURL
}

// New 创建一个新的 Token
func New(tokenType Type) *Token {
	return &Token{
		Type:      tokenType,
		authToken: NewAuthTokenInfo(),
	}
}

// BotToken 机器人身份的 token
func BotToken(appID uint64, clientSecret string) *Token {
	return &Token{
		appID:        appID,
		clientSecret: clientSecret,
		authToken:    NewAuthTokenInfo(),
		Type:         TypeBot,
	}
}

// UserToken 用户身份的token
func UserToken(appID uint64, clientSecret string) *Token {
	return &Token{
		appID:        appID,
		clientSecret: clientSecret,
		authToken:    NewAuthTokenInfo(),
		Type:         TypeNormal,
	}
}

// InitToken 启动机器人Token自动刷新
func (t *Token) InitToken(ctx context.Context) (err error) {
	if err = t.authToken.StartRefreshAccessToken(ctx, t.tokenURL, fmt.Sprint(t.appID), t.clientSecret); err != nil {
		return err
	}
	return nil
}

// GetAppID 取得Token中的appid
func (t *Token) GetAppID() uint64 {
	return t.appID
}

// GetString 获取授权头字符串
func (t *Token) GetString() string {
	if t.Type == TypeNormal || t.Type == TypeQQBot {
		return "QQBot " + t.GetAccessToken()
	}
	return fmt.Sprintf("%v.%s", t.appID, t.GetAccessToken())
}

// GetAccessToken 取得鉴权Token
func (t *Token) GetAccessToken() string {
	return t.authToken.getAuthToken().Token
}

// UpAccessToken 更新accessToken
func (t *Token) UpAccessToken(ctx context.Context, reason interface{}) error {
	return t.authToken.ForceUpToken(ctx, fmt.Sprint(reason))
}

// LoadFromConfig 从配置中读取 appid 和 token
func (t *Token) LoadFromConfig(file string) error {
	var conf struct {
		AppID uint64 `yaml:"appid"`
		Token string `yaml:"token"`
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Errorf("read token from file failed, err: %v", err)
		return err
	}
	if err = yaml.Unmarshal(content, &conf); err != nil {
		log.Errorf("parse config failed, err: %v", err)
		return err
	}
	t.appID = conf.AppID
	t.clientSecret = conf.Token
	return nil
}
