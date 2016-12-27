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
	Topics []string `json:"topics"`
}

type Unsubscribe struct {
	Topics []string `json:"topics"`
}
