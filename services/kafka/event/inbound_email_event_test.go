package event_test

import (
	"testing"

	"github.com/nirnanaaa/cloudive-mailer/services/kafka/event"
)

var expectedOutput = `{"trace_id":"abc","recipient":{"name":"","email":"","tracking_id":""},"sender":{"name":"","email":"","tracking_id":""},"subject":"","payload":null,"attachments":null}`

func TestS3_Encoding(t *testing.T) {
	// Parse configuration.
	c := event.InboundEmailEvent{TraceID: "abc"}
	s, err := event.EncodeOutgoingEvent(&c)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(s) != expectedOutput {
		t.Fatalf("Failed to match encoding output, got %s, expected %s", s, expectedOutput)
	}
}
func TestS3_Decoding(t *testing.T) {
	// Parse configuration.
	c := event.InboundEmailEvent{TraceID: "abc"}
	s, err := event.DecodeIncomingEvent([]byte(expectedOutput))
	if err != nil {
		t.Fatal(err.Error())
	}
	if s.TraceID != c.TraceID {
		t.Fatalf("Failed to match encoding output, got %v, expected %s", s, expectedOutput)
	}
}
