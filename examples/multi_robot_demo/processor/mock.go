package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

type mockProcessor struct {
	api openapi.OpenAPI
}

// GetReplayMsg 通用对话机器人处理逻辑
func (c *mockProcessor) GetReplayMsg(ctx context.Context, recMsg dto.Message) (dto.APIMessage,
	error) {
	return getTextMsg(recMsg.ID, recMsg.Content)
}

func (c *mockProcessor) SetAPI(api openapi.OpenAPI) {
	c.api = api
}

func (c *mockProcessor) GetAPI() openapi.OpenAPI {
	return c.api
}

func getTextMsg(id, msg string) (dto.MessageToCreate, error) {
	return dto.MessageToCreate{
		Timestamp: time.Now().UnixMilli(),
		Content:   fmt.Sprintf("mock processor:%s", msg),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             id,
			IgnoreGetMessageError: true,
		},
		MsgID: id,
	}, nil
}
