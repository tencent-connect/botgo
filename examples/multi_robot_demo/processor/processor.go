package processor

import (
	"context"
	"sync"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

// Processor ...
type Processor interface {
	GetAPI() openapi.OpenAPI
	SetAPI(openapi.OpenAPI)
	GetReplayMsg(ctx context.Context, recMsg dto.Message) (dto.APIMessage,
		error)
}

// NewProcessor ...
func NewProcessor(appID uint64) Processor {
	return &mockProcessor{}
}

var robotProc = make(map[uint64]Processor)
var lock sync.Mutex

// RegisterProcessor ...
func RegisterProcessor(appid uint64, proc Processor) {
	lock.Lock()
	defer lock.Unlock()
	robotProc[appid] = proc
}

// GetProcessor ...
func GetProcessor(appid uint64) Processor {
	return robotProc[appid]
}
