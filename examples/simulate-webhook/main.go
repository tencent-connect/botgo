package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/interaction/signature"
	"github.com/tencent-connect/botgo/interaction/webhook"
)

var handler event.PlainEventHandler = func(payload *dto.WSPayload, message []byte) error {
	fmt.Println(payload, message)
	return nil
}

const host = "http://localhost"
const port = ":8081"
const path = "/bot"
const url = host + port + path

func main() {
	event.RegisterHandlers(handler)
	http.HandleFunc(path, webhook.HTTPHandler)
	go simulateRequest()
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func simulateRequest() {
	// 等待 http 服务启动
	time.Sleep(3 * time.Second)
	var heartbeat = &dto.WSPayload{
		WSPayloadBase: dto.WSPayloadBase{
			OPCode: dto.WSHeartbeat,
		},
		Data: 123,
	}
	payload, _ := json.Marshal(heartbeat)
	send(payload)

	var dispatchEvent = &dto.WSPayload{
		WSPayloadBase: dto.WSPayloadBase{
			OPCode: dto.WSDispatchEvent,
			Seq:    1,
			Type:   dto.EventMessageReactionAdd,
		},
		Data: dto.WSMessageReactionData{
			UserID:    "123",
			ChannelID: "111",
			GuildID:   "222",
			Target: dto.ReactionTarget{
				ID:   "333",
				Type: dto.ReactionTargetTypeMsg,
			},
			Emoji: dto.Emoji{
				ID:   "42",
				Type: 1,
			},
		},
		RawMessage: nil,
	}
	payload, _ = json.Marshal(dispatchEvent)
	fmt.Println(string(payload))
	send(payload)
}

func send(payload []byte) {
	header := http.Header{}
	header.Set(signature.HeaderTimestamp, strconv.FormatUint(uint64(time.Now().Unix()), 10))

	sig, err := signature.Generate(webhook.DefaultGetSecretFunc(), header, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	header.Set(signature.HeaderSig, sig)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header = header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	r, _ := io.ReadAll(resp.Body)
	fmt.Printf("receive resp: %s", string(r))
}
