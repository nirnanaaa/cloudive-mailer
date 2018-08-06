package smtp

const (
	DefaultFromEmail              = "info@cloudive.cc"
	DefaultFromName               = "Cloudive"
	DefaultEnabled                = true
	DefaultSmtpPort               = 25
	DefaultHostName               = "mailcatcher"
	DefaultDomainWhitelistEnabled = false
)

// Config represents a configuration for a HTTP service.
type Config struct {
	Enabled                          bool     `toml:"enabled"`
	Hostname                         string   `toml:"hostname"`
	Port                             int      `toml:"port"`
	Username                         string   `toml:"username"`
	Password                         string   `toml:"password"`
	FromName                         string   `toml:"from-name"`
	FromMail                         string   `toml:"from-mail"`
	AttachmentDomainWhitelistEnabled bool     `toml:"domain-whitelist-enabled"`
	DomainWhitelist                  []string `toml:"domain-whitelist"`
}

// NewConfig returns a new Config with default settings.
func NewConfig() Config {
	return Config{
		AttachmentDomainWhitelistEnabled: DefaultDomainWhitelistEnabled,
		Enabled:         DefaultEnabled,
		Port:            DefaultSmtpPort,
		Hostname:        DefaultHostName,
		FromMail:        DefaultFromEmail,
		FromName:        DefaultFromName,
		DomainWhitelist: []string{},
	}
}
