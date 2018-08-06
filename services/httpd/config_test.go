package httpd_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/cloudive-mailer/services/httpd"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	c := httpd.NewConfig()
	if _, err := toml.Decode(`
		bind-address = "0.0.0.0:9009"
		enabled = true
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.BindAddress != "0.0.0.0:9009" {
		t.Fatalf("unexpected BindAddress: %s", c.BindAddress)
	} else if !c.Enabled {
		t.Fatalf("unexpected Enabled, want true got false")

	}
}
