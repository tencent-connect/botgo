// Package webhook HTTP回调处理
package webhook

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/tencent-connect/botgo/constant"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/interaction/signature"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/token"
)

type ack struct {
	Op   dto.OPCode `json:"op"`
	Data uint32     `json:"d"`
}

// GenHeartbeatACK 生成 http gateway 的心跳回包
func GenHeartbeatACK(seq uint32) string {
	s, _ := json.Marshal(ack{Op: dto.WSHeartbeatAck, Data: seq})
	return string(s)
}

// GenDispatchACK 生成事件包的回包，如果处理失败，则返回的 d 为 1，服务端会尝试重试
func GenDispatchACK(success bool) string {
	var r uint32
	if !success {
		r = 1
	}
	s, _ := json.Marshal(ack{Op: dto.HTTPCallbackAck, Data: r})
	return string(s)
}

// Deprecated: DefaultGetSecretFunc 默认的获取 secret 的函数，默认从环境变量读取
// 开发者如果需要从自己的配置文件，或者是其他地方获取 secret，可以重写这个函数
var DefaultGetSecretFunc = func() string {
	return os.Getenv("QQBotSecret")
}

// HTTPHandler 用户处理回调时间，该函数实现的是 https://pkg.go.dev/net/http#HandleFunc 所要求的 handler
// 会自动进行签名验证，心跳包回复，以及根据使用 event.RegisterHandlers 注册的 handler 去执行不同的 handler 来处理事件
// 如果开发者不想在接收事件的地方处理，可以实现 DefaultHandlers.Plain 然后在内部处理相关的异步生产或者转发的逻辑
func HTTPHandler(w http.ResponseWriter, r *http.Request, credentials *token.QQBotCredentials) {
	defer r.Body.Close()
	body := make([]byte, r.ContentLength)
	if _, err := r.Body.Read(body); err != nil && err != io.EOF {
		log.Errorf("read http callback body error: %s", err)
		return
	}
	log.Debugf("http callback body: %v", string(body))
	log.Debugf("http callback header: %v", r.Header)
	traceID := r.Header.Get(constant.HeaderTraceID)
	// 签名验证
	if pass, err := signature.Verify(credentials.AppSecret, r.Header, body); err != nil || !pass {
		log.Errorf("signature verify failed, err: %v, traceID: %s", err, traceID)
		return
	}
	// 解析 payload
	payload := &dto.WSPayload{}
	if err := json.Unmarshal(body, payload); err != nil {
		log.Errorf("unmarshal http callback body error: %s, traceID: %s", err, traceID)
		return
	}
	log.Info("payload:%+v", payload)
	// 原始数据放入，parse 的时候需要从里面提取 d
	payload.RawMessage = body
	payload.Session = &dto.Session{AppID: credentials.AppID}
	var result string
	if payload.OPCode == dto.HTTPCallbackValidation {
		data, ok := payload.Data.(map[string]interface{})
		if !ok {
			log.Errorf("callback data invalid: %+v, traceID: %s", payload.Data, traceID)
			return
		}
		plainToken, ptOk := data["plain_token"].(string)
		eventTs, etOk := data["event_ts"].(string)
		if !ptOk || !etOk {
			log.Errorf("callback data invalid: %+v, traceID: %s", payload.Data, traceID)
		}
		req := &dto.WHValidationReq{
			PlainToken: plainToken,
			EventTs:    eventTs,
		}
		validationRsp := GenValidationACK(req, r.Header, credentials.AppSecret)
		if validationRsp != nil {
			if _, err := w.Write(validationRsp); err != nil {
				log.Errorf("write http callback response error: %s, traceID: %s", err, traceID)
				return
			}
		}
		return
	}

	result = parsePayload(payload, traceID)
	if result != "" {
		if _, err := w.Write([]byte(result)); err != nil {
			log.Errorf("write http callback response error: %s, traceID: %s", err, traceID)
			return
		}
	}
}

func parsePayload(payload *dto.WSPayload, traceID string) string {
	// 处理心跳包
	if payload.OPCode == dto.WSHeartbeat {
		return GenHeartbeatACK(uint32(payload.Data.(float64)))
	}
	// 处理事件
	if payload.OPCode == dto.WSDispatchEvent {
		// 解析具体事件，并投递给业务注册的 handler
		if err := event.ParseAndHandle(payload); err != nil {
			log.Errorf(
				"parseAndHandle failed, %v, traceID:%s, payload: %v", err,
				traceID, payload,
			)
			return GenDispatchACK(false)
		}
		return GenDispatchACK(true)
	}

	return ""
}

// GenValidationACK 生成回调校验回包
func GenValidationACK(req *dto.WHValidationReq, header http.Header, secret string) []byte {
	h := header.Clone()
	h.Set(signature.HeaderTimestamp, req.EventTs)
	sig, err := signature.Generate(secret, h, []byte(req.PlainToken))
	if err != nil {
		log.Errorf("generate signature failed:%+v", err)
		return nil
	}
	rsp, err := json.Marshal(
		&dto.WHValidationRsp{
			PlainToken: req.PlainToken,
			Signature:  sig,
		})
	if err != nil {
		log.Errorf("handle validation failed:", err)
		return nil
	}
	return rsp
}
