package server

type Command struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
