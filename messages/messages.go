package messages

type Command struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Subscribe struct {
	Topic string `json:"topic"`
}

type Unsubscribe struct {
	Topic string `json:"topic"`
}
