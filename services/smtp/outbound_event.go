package smtp

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

// OutboundEmailEvent is used to fanin from kafka or HTTP.
type OutboundEmailEvent struct {
	// For tracing with zipkin or jaeger
	TraceID   string  `json:"trace_id,omitempty"`
	Recipient Contact `json:"recipient"`
	Sender    Contact `json:"sender"`

	Subject string `json:"subject"`
	Payload []byte `json:"payload"`

	// Just the download references. not sending high volume data.
	Attachments []Attachment `json:"attachments"`
}
