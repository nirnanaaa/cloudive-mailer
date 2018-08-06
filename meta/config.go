package meta

const (

	// DefaultLogLevel sets an optimistic log level
	DefaultLogLevel = "warn"
)

// Config represents a configuration for a Metrics service.
type Config struct {
	LogLevel string `toml:"log-level"`
}

// NewConfig returns a new Config with default settings.
func NewConfig() *Config {
	return &Config{
		LogLevel: DefaultLogLevel,
	}
}
