package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/interaction/signature"
	"github.com/tencent-connect/botgo/token"
	"gopkg.in/yaml.v2"
)

const host = "http://localhost"
const port = ":9000"
const path = "/qqbot"
const url = host + port + path

func main() {
	// 加载 appid 和 token
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln("load config file failed, err:", err)
	}
	credentials := &token.QQBotCredentials{}
	if err = yaml.Unmarshal(content, &credentials); err != nil {
		log.Fatalln("parse config failed, err:", err)
	}
	log.Println("credentials:", credentials)
	if err != nil {
		log.Fatalln(err)
	}
	go simulateRequest(credentials)
	var ln string
	fmt.Scanln()
	_, _ = fmt.Sscanln("%v", ln)
	fmt.Println("end")
}

func simulateRequest(credentials *token.QQBotCredentials) {
	// 等待 http 服务启动
	time.Sleep(3 * time.Second)
	var heartbeat = &dto.WSPayload{
		WSPayloadBase: dto.WSPayloadBase{
			OPCode: dto.WSHeartbeat,
		},
		Data: 123,
	}
	payload, _ := json.Marshal(heartbeat)
	send(payload, credentials)

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
	send(payload, credentials)
}

func send(payload []byte, credentials *token.QQBotCredentials) {
	header := http.Header{}
	header.Set(signature.HeaderTimestamp, strconv.FormatUint(uint64(time.Now().Unix()), 10))

	sig, err := signature.Generate(credentials.AppSecret, header, payload)
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
