package token

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tencent-connect/botgo/log"
)

// getAccessTokenURL 取得AccessToken的地址
var getAccessTokenURL = "https://bots.qq.com/app/getAppAccessToken"

// AuthTokenInfo 动态鉴权Token信息
type AuthTokenInfo struct {
	accessToken  AccessTokenInfo
	lock         *sync.RWMutex
	forceUpToken chan interface{}
	once         sync.Once
}

// AccessTokenInfo 鉴权Token信息
type AccessTokenInfo struct {
	Token     string
	ExpiresIn int64
	UpTime    time.Time
}

// NewAuthTokenInfo 初始化动态鉴权Token
func NewAuthTokenInfo() *AuthTokenInfo {
	return &AuthTokenInfo{
		lock:         &sync.RWMutex{},
		forceUpToken: make(chan interface{}, 10),
	}
}

// ForceUpToken 强制刷新Token
func (atoken *AuthTokenInfo) ForceUpToken(ctx context.Context, reason string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("uptoken err:%v", ctx.Err())
	case atoken.forceUpToken <- reason:
	}
	return nil
}

// StartRefreshAccessToken 启动获取AccessToken的后台刷新
func (atoken *AuthTokenInfo) StartRefreshAccessToken(ctx context.Context, tokenURL, appID, clientSecrent string) (err error) {
	tokenInfo, err := queryAccessToken(ctx, tokenURL, appID, clientSecrent)
	if err != nil {
		return err
	}
	atoken.setAuthToken(tokenInfo)
	tokenTTL := tokenInfo.ExpiresIn
	atoken.once.Do(func() {
		go func() {
			for {
				if tokenTTL <= 0 {
					tokenTTL = 1
				}
				select {
				case <-time.NewTimer(time.Duration(tokenTTL) * time.Second).C:
				case upToken := <-atoken.forceUpToken:
					log.Warnf("recv uptoken info:%v", upToken)
				case <-ctx.Done():
					log.Warnf("recv ctx:%v exit refreshAccessToken", ctx.Err())
					return
				}
				tokenInfo, err := queryAccessToken(ctx, tokenURL, appID, clientSecrent)
				if err == nil {
					atoken.setAuthToken(tokenInfo)
					tokenTTL = tokenInfo.ExpiresIn
				} else {
					log.Errorf("queryAccessToken err:%v", err)
				}
			}
		}()
	})
	return
}

func (atoken *AuthTokenInfo) getAuthToken() AccessTokenInfo {
	atoken.lock.RLock()
	defer atoken.lock.RUnlock()
	return atoken.accessToken

}

func (atoken *AuthTokenInfo) setAuthToken(accessToken AccessTokenInfo) {
	atoken.lock.Lock()
	defer atoken.lock.Unlock()
	atoken.accessToken = accessToken
}

type queryTokenReq struct {
	AppID        string `json:"appId"`
	ClientSecret string `json:"clientSecret"`
}

type queryTokenRsp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

func queryAccessToken(ctx context.Context, tokenURL, appID, clientSecrent string) (AccessTokenInfo, error) {
	method := "POST"

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("queryAccessToken err:%v", err)
		}
	}()
	if tokenURL == "" {
		tokenURL = getAccessTokenURL
	}

	queryReq := queryTokenReq{
		AppID:        appID,
		ClientSecret: clientSecrent,
	}
	data, err := json.Marshal(queryReq)
	if err != nil {
		return AccessTokenInfo{}, err
	}
	payload := bytes.NewReader(data)
	log.Infof("reqdata:%v", string(data))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(method, tokenURL, payload)
	if err != nil {
		log.Errorf("NewRequest err:%v", err)
		return AccessTokenInfo{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("http do err:%v", err)
		return AccessTokenInfo{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("ReadAll do err:%v", err)
		return AccessTokenInfo{}, err
	}
	log.Infof("accesstoken:%v", string(body))
	queryRsp := queryTokenRsp{}
	if err = json.Unmarshal(body, &queryRsp); err != nil {
		log.Errorf("Unmarshal err:%v", err)
		return AccessTokenInfo{}, err
	}

	rdata := AccessTokenInfo{
		Token:  queryRsp.AccessToken,
		UpTime: time.Now(),
	}
	rdata.ExpiresIn, _ = strconv.ParseInt(queryRsp.ExpiresIn, 10, 64)
	return rdata, err
}
