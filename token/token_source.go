// Package token 基于 golang.org/x/oauth2 标准实现token source
package token

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/tencent-connect/botgo/constant"
	"github.com/tencent-connect/botgo/log"
	"golang.org/x/oauth2"
	"golang.org/x/sync/singleflight"
)

const (
	// TypeBearer ..
	TypeBearer string = "Bearer"
	// TypeQQBot ..
	TypeQQBot string = "QQBot"

	defaultExpiryDeltaMillSec  = 9000 // 与oauth2.defaultExpiryDelta - time.Second
	randTimeUpperLimitMilliSec = 500  // 随机时间区间Sec
)

type qqBotTokenReq struct {
	AppID        string `json:"appId"`
	ClientSecret string `json:"clientSecret"`
}

type qqBotTokenRsp struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (r *qqBotTokenRsp) UnmarshalJSON(data []byte) error {
	// 创建一个临时结构体来解析 JSON 数据
	var temp struct {
		Code        int    `json:"code"`
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}

	// 解析 JSON 数据到临时结构体
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// 将字符串转换为 int64
	expiresIn, err := strconv.ParseInt(temp.ExpiresIn, 10, 64)
	if err != nil {
		return err
	}

	// 赋值给结构体字段
	r.ExpiresIn = expiresIn
	r.Code = temp.Code
	r.Message = temp.Message
	r.AccessToken = temp.AccessToken
	return nil
}

// QQBotCredentials QQ机器人appid、secret
type QQBotCredentials struct {
	AppID     string `yaml:"appid"`
	AppSecret string `yaml:"secret"`
}

// QQBotTokenSource QQ机器人token source
type QQBotTokenSource struct {
	credentials *QQBotCredentials
	cachedToken atomic.Value
	sg          singleflight.Group
}

// NewQQBotTokenSource 初始化
func NewQQBotTokenSource(credentials *QQBotCredentials) oauth2.TokenSource {
	return &QQBotTokenSource{
		credentials: credentials,
	}
}

// Token 获取access token
func (w *QQBotTokenSource) Token() (*oauth2.Token, error) {
	rawToken := w.cachedToken.Load()
	if rawToken != nil && rawToken.(*oauth2.Token).Valid() {
		token, ok := rawToken.(*oauth2.Token)
		if ok && token.Valid() {
			return token, nil
		}
	}
	// 获取新的access rawToken
	newToken, err, shard := w.sg.Do("retrieve access rawToken", func() (interface{}, error) {
		return w.getNewToken()
	})
	log.Debugf("shared flight:%v", shard)
	if err != nil {
		return nil, err
	}
	w.cachedToken.Store(newToken)
	return newToken.(*oauth2.Token), nil
}

func (w *QQBotTokenSource) getNewToken() (*oauth2.Token, error) {
	retrieveReq := qqBotTokenReq{
		AppID:        w.credentials.AppID,
		ClientSecret: w.credentials.AppSecret,
	}
	data, err := json.Marshal(retrieveReq)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewReader(data)
	log.Debugf("retrieve access token URL:%v req:%v", getTokenURL(), string(data))
	req, err := http.NewRequest(http.MethodPost, getTokenURL(), payload)
	if err != nil {
		log.Errorf("init http req failed:%v", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	rsp, err := client.Do(req)
	if err != nil {
		log.Errorf("retrieve token failed:%v", err)
		return nil, err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	rspTraceID := rsp.Header.Get(constant.HeaderTraceID)
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Errorf("read rsp failed:%v", err)
		return nil, err
	}
	log.Debugf("access token:%v traceID:%v", string(body), rspTraceID)
	retrieveRsp := &qqBotTokenRsp{}
	if err = json.Unmarshal(body, retrieveRsp); err != nil {
		log.Errorf("unmarshal rsp failed:%v traceID:%v", err, rspTraceID)
		return nil, err
	}
	if retrieveRsp.Code != 0 {
		log.Errorf("query acessToken err:%v.%v traceID:%v", retrieveRsp.Code, retrieveRsp.Message, rspTraceID)
		return nil, fmt.Errorf("%v.%v", retrieveRsp.Code, retrieveRsp.Message)
	}
	expiry := time.Now().Add(time.Duration(retrieveRsp.ExpiresIn) * time.Second)
	return &oauth2.Token{
		AccessToken: retrieveRsp.AccessToken,
		TokenType:   TypeQQBot,
		Expiry:      expiry,
		ExpiresIn:   retrieveRsp.ExpiresIn,
	}, nil
}

// GetAppID 获取appid
func (w *QQBotTokenSource) GetAppID() string {
	if w == nil || w.credentials == nil {
		return ""
	}
	return w.credentials.AppID
}

// StartRefreshAccessToken 启动获取AccessToken的后台刷新
func StartRefreshAccessToken(ctx context.Context, tokenSource oauth2.TokenSource) error {
	tk, err := tokenSource.Token()
	if err != nil {
		return err
	}
	log.Debugf("token:%+v ", tk)
	go func() {
		var consecutiveFailures int
		for {
			var refreshMilliSec int64
			//上一轮获取 tk 失败
			if tk == nil {
				if consecutiveFailures > 10 {
					panic("get token failed continuously for more than ten times")
				}
				consecutiveFailures++
				refreshMilliSec = 1000 // 1000ms后重试
			} else {
				consecutiveFailures = 0
				refreshMilliSec = getRefreshMilliSec(tk.ExpiresIn)
			}
			log.Debugf("refresh after %d milli sec", refreshMilliSec)
			timer := time.NewTimer(time.Duration(refreshMilliSec) * time.Millisecond)
			select {
			case <-timer.C:
				{
					log.Debugf("start to refresh access token %s", time.Now().Format(time.StampMilli))
					tk, err = tokenSource.Token()
					if err != nil {
						log.Errorf("refresh access token failed:%s", err)
					}
					timer.Stop()
				}
			case <-ctx.Done():
				{
					log.Warnf("recv ctx:%v exit refresh token", ctx.Err())
					timer.Stop()
					return
				}
			}
		}
	}()
	return nil
}

var (
	r = rand.New(rand.NewSource(time.Now().Unix()))
)

// getRefreshSec 为token刷新保留提前量。避免由于网络延迟等原因导致的token刷新不及时。
func getRefreshMilliSec(tokenTTLSec int64) int64 {
	refreshMilliSec := tokenTTLSec * 1000
	if refreshMilliSec < defaultExpiryDeltaMillSec {
		return refreshMilliSec
	}
	refreshMilliSec -= defaultExpiryDeltaMillSec
	// 随机化，避免所有机器人都同时获取access_token
	if refreshMilliSec > randTimeUpperLimitMilliSec {
		rand := r.Int63n(randTimeUpperLimitMilliSec)
		log.Debugf("rand:%d", rand)
		refreshMilliSec -= rand
	}
	return refreshMilliSec
}

func getTokenURL() string {
	return fmt.Sprintf("%v%v", constant.TokenDomain, "/app/getAppAccessToken")
}
