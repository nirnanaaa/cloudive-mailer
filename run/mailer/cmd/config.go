package run

import (
	"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	itoml "github.com/influxdata/influxdb/toml"
	"github.com/nirnanaaa/cloudive-mailer/meta"
	"github.com/nirnanaaa/cloudive-mailer/services/httpd"
	"github.com/nirnanaaa/cloudive-mailer/services/kafka"
	"github.com/nirnanaaa/cloudive-mailer/services/smtp"
)

// Config represents the configuration format for the roove binary.
type Config struct {
	Kafka *kafka.Config `toml:"kafka"`
	Meta  *meta.Config  `toml:"meta"`
	HTTPD *httpd.Config `toml:"httpd"`
	SMTP  smtp.Config   `toml:"smtp"`
}

// NewConfig returns an instance of Config with reasonable defaults.
func NewConfig() *Config {
	c := &Config{}
	c.Kafka = kafka.NewConfig()
	c.Meta = meta.NewConfig()
	c.HTTPD = httpd.NewConfig()
	c.SMTP = smtp.NewConfig()
	return c
}

// NewDemoConfig returns the config that runs when no config is specified.
func NewDemoConfig() (*Config, error) {
	c := NewConfig()

	return c, nil
}

func trimBOM(f []byte) []byte {
	return bytes.TrimPrefix(f, []byte("\xef\xbb\xbf"))
}

// FromTomlFile loads the config from a TOML file.
func (c *Config) FromTomlFile(fpath string) error {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	bs = trimBOM(bs)
	return c.FromToml(string(bs))
}

// FromToml loads the config from TOML.
func (c *Config) FromToml(input string) error {
	// Replace deprecated [cluster] with [coordinator]
	re := regexp.MustCompile(`(?m)^\s*\[cluster\]`)
	input = re.ReplaceAllStringFunc(input, func(in string) string {
		in = strings.TrimSpace(in)
		out := "[coordinator]"
		log.Printf("deprecated config option %s replaced with %s; %s will not be supported in a future release\n", in, out, in)
		return out
	})

	_, err := toml.Decode(input, c)
	return err
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	return nil
}

// ApplyEnvOverrides apply the environment configuration on top of the config.
func (c *Config) ApplyEnvOverrides(getenv func(string) string) error {
	return itoml.ApplyEnvOverrides(getenv, "ROOVE", c)
}
