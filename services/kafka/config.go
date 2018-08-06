package kafka

const (
	// DefaultBroker defines the default broker for kafka
	DefaultBroker = "localhost:9092"

	// DefaultInboundQueue defines the name of the inbound queue
	DefaultInboundQueue = "mail"

	// DefaultOutboundQueue defines an outbound queue name
	DefaultOutboundQueue = "mail-worker-queue"

	// DefaultGroupName defines a group name for processing
	DefaultGroupName = "mail-processor-name"
)

// Config represents a configuration for a Kafka service.
type Config struct {
	Brokers           []string `toml:"brokers"`
	InboundQueueName  string   `toml:"inbound-queue"`
	OutboundQueueName string   `toml:"outbound-queue"`
	GroupName         string   `toml:"group"`
}

// NewConfig returns a new Config with default settings.
func NewConfig() *Config {
	return &Config{
		Brokers:           []string{DefaultBroker},
		InboundQueueName:  DefaultInboundQueue,
		OutboundQueueName: DefaultOutboundQueue,
		GroupName:         DefaultGroupName,
	}
}
