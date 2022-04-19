package search

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/interaction/signature"
	"github.com/tencent-connect/botgo/log"
)

const maxRespBuffer = 65535

// Config 搜索请求配置
type Config struct {
	AppID    string
	EndPoint string // 回调url地址
	Secret   string
}

// SimulateSearch 模拟内联搜索请求
// 开发者可以使用本方法请求自己的服务器进行平台内联搜索的模拟，避免在平台上触发搜索请求。提升联调效率。
func SimulateSearch(config *Config, keyword string) (*dto.SearchRsp, error) {
	interactionData := &dto.InteractionData{
		Name: "search",
		Type: dto.InteractionDataTypeChatSearch,
	}
	interactionData.Resolved, _ = json.Marshal(dto.SearchInputResolved{Keyword: keyword})
	interaction := &dto.Interaction{
		ApplicationID: config.AppID,
		Type:          dto.InteractionTypeCommand,
		Data:          interactionData,
		Version:       1,
	}
	jsonStr, _ := json.Marshal(interaction)
	timestamp := strconv.FormatUint(uint64(time.Now().Unix()), 10)

	// calc sig
	header := http.Header{}
	header.Set(signature.HeaderTimestamp, timestamp)
	sig, err := signature.Generate(config.Secret, header, jsonStr)
	if err != nil {
		return nil, err
	}

	// build req
	req, err := http.NewRequest(http.MethodPost, config.EndPoint, bytes.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set(signature.HeaderTimestamp, timestamp)
	req.Header.Set(signature.HeaderSig, sig)
	log.Info(req)

	// parse resp
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Info(resp)
	defer func() {
		resp.Body.Close()
	}()

	// parse resp body
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, maxRespBuffer))
	if err != nil {
		return nil, err
	}
	log.Info(string(body))
	result := &dto.SearchRsp{}
	if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil
}
