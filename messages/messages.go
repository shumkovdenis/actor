package messages

type Command struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Subscribe -> command.subscribe
type Subscribe struct {
	Topics []string `mapstructure:"topics"`
}

// SubscribeSuccess -> event.subscribe.success
type SubscribeSuccess struct {
	Topics []string `json:"topics"`
}

func (s *SubscribeSuccess) Contains(topic string) bool {
	for _, t := range s.Topics {
		if topic == t {
			return true
		}
	}
	return false
}

// Unsubscribe -> command.unsubscribe
type Unsubscribe struct {
	Topics []string `mapstructure:"topics"`
}

// UnsubscribeSuccess -> event.unsubscribe.success
type UnsubscribeSuccess struct {
	Topics []string `json:"topics"`
}

func (s *UnsubscribeSuccess) Contains(topic string) bool {
	for _, t := range s.Topics {
		if topic == t {
			return true
		}
	}
	return false
}
