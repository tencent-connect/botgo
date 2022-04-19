package dto

import "encoding/json"

// Interaction 互动行为对象
type Interaction struct {
	ID            string           `json:"id,omitempty"`             // 互动行为唯一标识
	ApplicationID string           `json:"application_id,omitempty"` // 应用ID
	Type          InteractionType  `json:"type,omitempty"`           // 互动类型
	Data          *InteractionData `json:"data,omitempty"`           // 互动数据
	GuildID       string           `json:"guild_id,omitempty"`       // 频道 ID
	ChannelID     string           `json:"channel_id,omitempty"`     // 子频道 ID
	Version       uint32           `json:"version,omitempty"`        //	版本，默认为 1
}

// InteractionType 互动类型
type InteractionType uint32

const (
	// InteractionTypePing ping
	InteractionTypePing InteractionType = 1
	// InteractionTypeCommand 命令
	InteractionTypeCommand InteractionType = 2
)

// InteractionData 互动数据
type InteractionData struct {
	Name     string              `json:"name,omitempty"`     // 标题
	Type     InteractionDataType `json:"type,omitempty"`     //	数据类型，不同数据类型对应不同的 resolved 数据
	Resolved json.RawMessage     `json:"resolved,omitempty"` // 跟不同的互动类型和数据类型有关系的数据
}

// InteractionDataType 互动数据类型
type InteractionDataType uint32

const (
	// InteractionDataTypeChatSearch 聊天框搜索
	InteractionDataTypeChatSearch InteractionDataType = 9
)
