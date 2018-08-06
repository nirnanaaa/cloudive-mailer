package httpd

const (
	// DefaultBindAddress defines a bind address for the http service
	DefaultBindAddress = "127.0.0.1:9009"

	// DefaultEnabled enables or disables the HTTP service
	DefaultEnabled = false
)

// Config represents a configuration for a Kafka service.
type Config struct {
	Enabled     bool   `toml:"enabled"`
	BindAddress string `toml:"bind-address"`
}

// NewConfig returns a new Config with default settings.
func NewConfig() *Config {
	return &Config{
		Enabled:     DefaultEnabled,
		BindAddress: DefaultBindAddress,
	}
}
