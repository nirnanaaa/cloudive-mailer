package run_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	run "github.com/nirnanaaa/cloudive-mailer/run/mailer/cmd"
	"github.com/sirupsen/logrus"
)

// Ensure the configuration can be parsed.
func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	var c run.Config
	if err := c.FromToml(`
		[meta]
		log-level = "debug"
	  
	  [httpd]
		enabled = true
		bind-address = ":8087"
	  
	  [kafka]
		brokers = ["localhost:9092"]
		inbound-queue = "s3notifications"
		outbound-queue = "thumb-worker-queue"
	  
	  [s3]
		ssl-enabled = false
		endpoint = "localhost:9000"
		access-key-id = "abc"
		secret-access-key = "abc"
`); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.Meta.LogLevel != "debug" {
		t.Fatalf("unexpected log level: %s", c.Meta.LogLevel)
	} else if !c.HTTPD.Enabled {
		t.Fatalf("unexpected http enabled: false")
	} else if c.HTTPD.BindAddress != ":8087" {
		t.Fatalf("unexpected api bind address: %s", c.HTTPD.BindAddress)
	}
}

// Ensure the configuration can be parsed.
func TestConfig_Parse_EnvOverride(t *testing.T) {
	// Parse configuration.
	c := run.NewConfig()
	if _, err := toml.Decode(`
		[meta]
		log-level = "debug"
	  
	  [httpd]
		enabled = true
		bind-address = ":8087"
	  
	  [kafka]
		brokers = ["localhost:9092"]
		inbound-queue = "s3notifications"
		outbound-queue = "thumb-worker-queue"
	  
	  [s3]
		ssl-enabled = false
		endpoint = "localhost:9000"
		access-key-id = "abc"
		secret-access-key = "abc"
`, &c); err != nil {
		t.Fatal(err)
	}

	getenv := func(s string) string {
		switch s {
		case "CLOUDIVE_KAFKA_BROKERS":
			return "kafka-1:9092"
		case "CLOUDIVE_HTTPD_BIND_ADDRESS":
			return ":1234"
		case "CLOUDIVE_META_LOG_LEVEL":
			return "warning"
		}
		return ""
	}

	if err := c.ApplyEnvOverrides(getenv); err != nil {
		t.Fatalf("failed to apply env overrides: %v", err)
	}
	if c.Kafka.Brokers[0] != "kafka-1:9092" {
		t.Fatalf("unexpected kafka broker list: %+v", c.Kafka.Brokers)
	}
	if c.HTTPD.BindAddress != ":1234" {
		t.Fatalf("unexpected httpd bind address: %s", c.HTTPD.BindAddress)
	}

	if c.Meta.LogLevel != logrus.WarnLevel.String() {
		t.Fatalf("unexpected logging level: %v", c.Meta.LogLevel)
	}
}
