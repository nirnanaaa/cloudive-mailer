package kafka

import (
	"encoding/json"
	"time"
)

// Meta file info
type Meta struct {
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
}

// Event for kafka outputs
type Event struct {
	Key         string    `json:"key"`
	Meta        Meta      `json:"meta"`
	ProcessedAt time.Time `json:"processedAt"`
}

// EncodeOutgoingEvent encodes an outgoing kafka event
func EncodeOutgoingEvent(evt *InboundEmailEvent) ([]byte, error) {
	data, err := json.Marshal(evt)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// DecodeIncomingEvent decodes an incoming kafka event
func DecodeIncomingEvent(message []byte) (*InboundEmailEvent, error) {
	var ins InboundEmailEvent
	if err := json.Unmarshal(message, &ins); err != nil {
		return nil, err
	}
	return &ins, nil
}
