package event

import "encoding/json"

type Attachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Contact defines either a recipient or sender
type Contact struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	TrackingID string `json:"tracking_id"`
}

// InboundEmailEvent is used to fanin from kafka or HTTP.
type InboundEmailEvent struct {
	// For tracing with zipkin or jaeger
	TraceID   string  `json:"trace_id,omitempty"`
	Recipient Contact `json:"recipient"`
	Sender    Contact `json:"sender"`

	Subject string `json:"subject"`
	Payload []byte `json:"payload"`

	// Just the download references. not sending high volume data.
	Attachments []Attachment `json:"attachments"`
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
