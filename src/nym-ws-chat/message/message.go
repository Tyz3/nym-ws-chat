package message

import "encoding/json"

type Message interface {
	ToJson() []byte
}

type OneWayMessage struct {
	Type          string `json:"type"`
	Recipient     string `json:"recipient"`
	Message       string `json:"message"`
	WithReplySurb bool   `json:"withReplySurb"`
}

func NewOneWayMessage(text string, recipient string, surb bool) *OneWayMessage {
	return &OneWayMessage{
		Type:          "send",
		Message:       text,
		Recipient:     recipient,
		WithReplySurb: surb,
	}
}

func (m *OneWayMessage) ToJson() []byte {
	data, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}

	return data
}
