package dto

// HTTPIdentity 鉴权数据
type HTTPIdentity struct {
	Intents  Intent    `json:"intents"`
	Shards   [2]uint32 `json:"shards"` // array of two integers (shard_id, num_shards)
	Callback string    `json:"callback_url"`
}

// HTTPReady ready，鉴权后返回
type HTTPReady struct {
	Version   int    `json:"version"`
	SessionID string `json:"session_id"`
	Bot       struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"bot"`
	Shard [2]uint32 `json:"shard"`
}

// HTTPSession session 对象
type HTTPSession struct {
	AppID             int64    `json:"app_id"`
	SessionID         string   `json:"session_id"`
	CallbackURL       string   `json:"callback_url"`
	Env               string   `json:"env"`
	Intents           int64    `json:"intents"`
	LastHeartbeatTime string   `json:"last_heartbeat_time"`
	State             string   `json:"state"`
	Shards            [2]int64 `json:"shards"`
}
