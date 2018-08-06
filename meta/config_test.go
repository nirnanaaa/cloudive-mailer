package meta_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/cloudive-mailer/meta"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	c := meta.NewConfig()
	if _, err := toml.Decode(`
		log-level = "warn"
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.LogLevel != "warn" {
		t.Fatalf("unexpected log level: %s", c.LogLevel)
	}
}
