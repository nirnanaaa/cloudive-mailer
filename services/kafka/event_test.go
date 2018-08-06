package kafka_test

import (
	"testing"

	"github.com/nirnanaaa/cloudive-mailer/services/kafka"
)

var expectedOutput = `{"key":"abc","s3Event":{"eventVersion":"","eventSource":"","awsRegion":"","eventTime":"","eventName":"s3:ObjectAccessed:Get","userIdentity":{"principalId":""},"requestParameters":null,"responseElements":null,"s3":{"s3SchemaVersion":"","configurationId":"","bucket":{"name":"","ownerIdentity":{"principalId":""},"arn":""},"object":{"key":"","sequencer":""}},"source":{"host":"","port":"","userAgent":""}},"format":{"Width":0,"Height":0,"Quality":0,"Name":""},"meta":{"size":0,"mimeType":""},"processedAt":"0001-01-01T00:00:00Z"}`

func TestS3_Encoding(t *testing.T) {
	// Parse configuration.
	c := kafka.InboundEmailEvent{TraceID: "abc"}
	s, err := kafka.EncodeOutgoingEvent(&c)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(s) != expectedOutput {
		t.Fatalf("Failed to match encoding output, got %s, expected %s", s, expectedOutput)
	}
}
func TestS3_Decoding(t *testing.T) {
	// Parse configuration.
	c := kafka.InboundEmailEvent{TraceID: "abc"}
	s, err := kafka.DecodeIncomingEvent([]byte(expectedOutput))
	if err != nil {
		t.Fatal(err.Error())
	}
	if s.TraceID != c.TraceID {
		t.Fatalf("Failed to match encoding output, got %v, expected %s", s, expectedOutput)
	}
}
