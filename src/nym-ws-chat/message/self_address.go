package message

import "encoding/json"

type SelfAddressMessage struct {
	Type string `json:"type"`
}

func NewSelfAddressMessage() *SelfAddressMessage {
	return &SelfAddressMessage{
		Type: "selfAddress",
	}
}

func (m *SelfAddressMessage) ToJson() []byte {
	data, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}

	return data
}
