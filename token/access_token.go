package token

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/tencent-connect/botgo/constant"
	"github.com/tencent-connect/botgo/log"
)

var (
	r *rand.Rand
)

const (
	preserveTokenTTL   float64 = 30 // token预留时长，用于控制提前刷新token Sec
	minTimeGap         float64 = 2  // 定时器的最少时长Sec
	randTimeUpperLimit         = 10 // 随机时间区间Sec

	tokenURL = "/app/getAppAccessToken" // tokenURL 取得AccessToken的地址

)

func init() {
	src := rand.NewSource(time.Now().Unix())
	r = rand.New(src)
}

// AccessToken 鉴权Token信息
type AccessToken struct {
	Token      string
	ExpiresIn  int64
	UpdateTime time.Time
}

// GetToken 获取 token
func (a *AccessToken) GetToken() string {
	if a == nil {
		return ""
	}
	return a.Token
}

type retrieveTokenReq struct {
	AppID        string `json:"appId"`
	ClientSecret string `json:"clientSecret"`
}

type retrieveTokenRsp struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

func getTokenURL() string {
	return fmt.Sprintf("%v%v", constant.TokenDomain, tokenURL)
}

func retrieveToken(appID, secret string) (*AccessToken, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("query access token err:%v", err)
		}
	}()
	retrieveReq := retrieveTokenReq{
		AppID:        appID,
		ClientSecret: secret,
	}
	data, err := json.Marshal(retrieveReq)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewReader(data)
	log.Debugf("retrieve access token URL:%v req:%v", getTokenURL(), string(data))

	req, err := http.NewRequest(http.MethodPost, getTokenURL(), payload)
	if err != nil {
		log.Errorf("init http request failed:%v", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("retrieve token failed:%v", err)
		return nil, err
	}
	defer res.Body.Close()

	rsptraceID := res.Header.Get(constant.TraceIDKey)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("read rsp failed:%v", err)
		return nil, err
	}
	log.Debugf("access token:%v traceID:%v", string(body), rsptraceID)
	retrieveRsp := &retrieveTokenRsp{}
	if err = json.Unmarshal(body, retrieveRsp); err != nil {
		log.Errorf("unmarshal rsp failed:%v traceID:%v", err, rsptraceID)
		return nil, err
	}
	if retrieveRsp.Code != 0 {
		log.Errorf("query acessToken err:%v.%v traceID:%v", retrieveRsp.Code, retrieveRsp.Message, rsptraceID)
		return nil, fmt.Errorf("%v.%v", retrieveRsp.Code, retrieveRsp.Message)
	}
	rdata := &AccessToken{
		Token:      retrieveRsp.AccessToken,
		UpdateTime: time.Now(),
	}
	rdata.ExpiresIn, err = strconv.ParseInt(retrieveRsp.ExpiresIn, 10, 64)
	if err != nil {
		log.Errorf("parse expire_in failed err:%v traceID:%v", err, rsptraceID)
		return nil, err
	}
	return rdata, nil
}

// getTokenTTL 为token刷新保留提前量。避免由于网络延迟等原因导致的token刷新不及时。
func getTokenTTL(tokenTTL float64) float64 {
	if tokenTTL < preserveTokenTTL {
		return tokenTTL
	}
	tokenTTL = tokenTTL - preserveTokenTTL
	if tokenTTL < minTimeGap {
		return minTimeGap // 为了避免有bug导致不断触发timer，这里需要预留一点时间
	}
	// 随机化，避免所有机器人都同时获取access_token
	if tokenTTL > randTimeUpperLimit {
		tokenTTL = tokenTTL - float64(r.Int63n(randTimeUpperLimit))
	}
	return tokenTTL
}
