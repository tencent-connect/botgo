package dto

// MessageArk ark模板消息
type MessageArk struct {
	Ark *Ark `json:"ark,omitempty"`
}

// Ark 消息模版
type Ark struct {
	TemplateID int      `json:"template_id,omitempty"` // ark 模版 ID
	KV         []*ArkKV `json:"kv,omitempty"`          // ArkKV 数组
}

// ArkKV Ark 键值对
type ArkKV struct {
	Key   string    `json:"key,omitempty"`
	Value string    `json:"value,omitempty"`
	Obj   []*ArkObj `json:"obj,omitempty"`
}

// ArkObj Ark 对象
type ArkObj struct {
	ObjKV []*ArkObjKV `json:"obj_kv,omitempty"`
}

// ArkObjKV Ark 对象键值对
type ArkObjKV struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
