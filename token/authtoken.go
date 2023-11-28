package token

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tencent-connect/botgo/log"
)

var (
	r *rand.Rand
)

func init() {
	src := rand.NewSource(time.Now().Unix())
	r = rand.New(src)
}

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

// NewAuthTokenInfo 创建 authToken 信息
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
func (atoken *AuthTokenInfo) StartRefreshAccessToken(ctx context.Context,
	tokenURL, appID, clientSecrent string) (err error) {
	tokenInfo, err := queryAccessToken(ctx, tokenURL, appID, clientSecrent)
	if err != nil {
		return err
	}
	atoken.setAuthToken(tokenInfo)
	tokenTTL := tokenInfo.ExpiresIn
	atoken.once.Do(func() {
		go func() {
			for {
				tokenTTL = getTokenTTL(tokenTTL)
				log.Infof("tokenTTL:%d", tokenTTL)
				select {
				case <-time.NewTimer(time.Duration(tokenTTL) * time.Second).C:
				case reason := <-atoken.forceUpToken:
					log.Warnf("forceUpToken, reason:%v", reason)
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

const (
	preserveTokenTTL = 30 // token预留时长，用于控制提前刷新token
	minTimer         = 2  // 定时器的最少时长
	randTime         = 10 // 随机时间区间
)

func getTokenTTL(tokenTTL int64) int64 {
	tokenTTL = tokenTTL - preserveTokenTTL
	if tokenTTL <= minTimer {
		tokenTTL = minTimer // 为了避免有bug导致不断触发timer，这里需要预留一点时间
	}
	// 随机化，避免所有机器人都同时获取access_token
	if tokenTTL > randTime {
		tokenTTL = tokenTTL - r.Int63n(randTime)
	}
	return tokenTTL
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
	log.Infof("queryAccessToken reqData:%v", string(data))
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
	log.Infof("access_token:%v", string(body))
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
