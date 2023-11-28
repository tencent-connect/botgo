package token

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"encoding/json"
	"github.com/tencent-connect/botgo/log"
	"math/rand"
	"net/http"
)

var (
	r *rand.Rand
)

const (
	preserveTokenTTL   float64 = 30                                          // token预留时长，用于控制提前刷新token Sec
	minTimeGap         float64 = 2                                           // 定时器的最少时长Sec
	randTimeUpperLimit         = 10                                          // 随机时间区间Sec
	tokenURL                   = "https://bots.qq.com/app/getAppAccessToken" // tokenURL 取得AccessToken的地址
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
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
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
	log.Debugf("retrieve access token req:%v", string(data))

	req, err := http.NewRequest(http.MethodPost, tokenURL, payload)
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("read rsp failed:%v", err)
		return nil, err
	}
	log.Debugf("access token:%v", string(body))
	retrieveRsp := &retrieveTokenRsp{}
	if err = json.Unmarshal(body, retrieveRsp); err != nil {
		log.Errorf("unmarshal rsp failed:%v", err)
		return nil, err
	}
	rdata := &AccessToken{
		Token:      retrieveRsp.AccessToken,
		UpdateTime: time.Now(),
	}
	rdata.ExpiresIn, err = strconv.ParseInt(retrieveRsp.ExpiresIn, 10, 64)
	if err != nil {
		log.Errorf("parse expire_in failed err:%v", err)
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
