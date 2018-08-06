package kafka_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/cloudive-mailer/services/kafka"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	c := kafka.NewConfig()
	if _, err := toml.Decode(`
		brokers = ["localhost:9092"]
		inbound-queue = "s3notifications"
		outbound-queue = "thumb-worker-queue"
		group = "s3-brokers"
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if len(c.Brokers) != 1 {
		t.Fatalf("unexpected broker count: %d", len(c.Brokers))
	} else if c.InboundQueueName != "s3notifications" {
		t.Fatalf("unexpected inbound queue name: %s", c.InboundQueueName)
	} else if c.OutboundQueueName != "thumb-worker-queue" {
		t.Fatalf("unexpected outbound queue name: %s", c.OutboundQueueName)
	} else if c.GroupName != "s3-brokers" {
		t.Fatalf("unexpected group name: %s", c.GroupName)
	}
}
