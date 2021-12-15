// Package token 用于调用 openapi，websocket 的 token 对象。
package token

import (
	"fmt"
	"io/ioutil"

	"github.com/tencent-connect/botgo/log"
	"gopkg.in/yaml.v3"
)

// Type token 类型
type Type string

// TokenType
const (
	TypeBot    Type = "Bot"
	TypeNormal Type = "Bearer"
)

// Token 用于调用接口的 token 结构
type Token struct {
	AppID       uint64
	AccessToken string
	Type        Type
}

// New 创建一个新的 Token
func New(tokenType Type) *Token {
	return &Token{
		Type: tokenType,
	}
}

// BotToken 机器人身份的 token
func BotToken(appID uint64, accessToken string) *Token {
	return &Token{
		AppID:       appID,
		AccessToken: accessToken,
		Type:        TypeBot,
	}
}

// UserToken 用户身份的token
func UserToken(appID uint64, accessToken string) *Token {
	return &Token{
		AppID:       appID,
		AccessToken: accessToken,
		Type:        TypeNormal,
	}
}

// GetString 获取授权头字符串
func (t *Token) GetString() string {
	if t.Type == TypeNormal {
		return t.AccessToken
	}
	return fmt.Sprintf("%v.%s", t.AppID, t.AccessToken)
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
	t.AppID = conf.AppID
	t.AccessToken = conf.Token
	return nil
}
