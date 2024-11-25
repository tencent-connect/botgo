package dto

// WHValidationReq 机器人回调验证请求Data
type WHValidationReq struct {
	PlainToken string `json:"plain_token"`
	EventTs    string `json:"event_ts"`
}

// WHValidationRsp 机器人回调验证响应结果
type WHValidationRsp struct {
	PlainToken  string `json:"plain_token"`
	Signature   string `json:"signature"`
	DataVersion string `json:"data_version"` //数据格式版本号
}
