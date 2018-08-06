package run

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/nirnanaaa/cloudive-mailer/services/httpd"
	"github.com/nirnanaaa/cloudive-mailer/services/kafka"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Service is an interface definition for each service.
type Service interface {
	Start() error
	Stop() error
}

// Command represents the command executed by "cloudive-thumber run".
type Command struct {
	Version   string
	Branch    string
	Commit    string
	BuildTime string

	closing chan struct{}
	err     chan error
	Closed  chan struct{}

	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Services []Service

	Logger *logrus.Logger
	Getenv func(string) string
}

// NewCommand return a new instance of Command.
func NewCommand() *Command {
	logger := log.New()
	logger.Formatter = new(prefixed.TextFormatter)
	return &Command{
		closing: make(chan struct{}),
		err:     make(chan error),
		Closed:  make(chan struct{}),
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Logger:  logger,
	}
}

const logo = `

`

// Run parses the config from args and runs the server.
func (cmd *Command) Run(args ...string) error {
	options, err := cmd.ParseFlags(args...)
	if err != nil {
		return err
	}
	logger := cmd.Logger
	// Print sweet InfluxDB logo.
	fmt.Print(logo)

	// Set parallelism.
	runtime.GOMAXPROCS(runtime.NumCPU())
	coreLog := logger.WithField("prefix", "core")
	// Mark start-up in log.
	coreLog.Infof("Parser starting, version %s, branch %s, commit %s",
		cmd.Version, cmd.Branch, cmd.Commit)
	coreLog.Infof("Go version %s, GOMAXPROCS set to %d", runtime.Version(), runtime.GOMAXPROCS(0))
	runtime.SetBlockProfileRate(int(1 * time.Second))
	config, err := cmd.ParseConfig(options.GetConfigPath())
	if err != nil {
		return fmt.Errorf("parse config: %s", err)
	}
	// Apply any environment variables on top of the parsed config
	if err := config.ApplyEnvOverrides(cmd.Getenv); err != nil {
		return fmt.Errorf("apply env config: %v", err)
	}

	// Validate the configuration.
	if err := config.Validate(); err != nil {
		return fmt.Errorf("%s. To generate a valid configuration file run `cloudive-thumber config > cloudive-thumber.generated.conf`", err)
	}
	level, err := log.ParseLevel(config.Meta.LogLevel)
	if err != nil {
		return err
	}
	logger.SetLevel(level)
	flag.Parse()
	kafkaService := kafka.NewService(config.Kafka)
	kafkaService.SetLogOutput(logger, "[kafka]")
	httpdService := httpd.NewService(*config.HTTPD)
	httpdService.SetLogOutput(logger)
	httpdService.Handler.Kafka = kafkaService
	// kafkaService.SetDefaultMessageProcessor()
	cmd.Services = append(cmd.Services, httpdService)
	cmd.Services = append(cmd.Services, kafkaService)
	return nil
}

// Open opens all services
func (cmd *Command) Open() error {
	for _, service := range cmd.Services {
		if err := service.Start(); err != nil {
			log.Errorf(err.Error())
			return fmt.Errorf("open service: %s", err)
		}
	}
	return nil
}

// Close shuts down the server.
func (cmd *Command) Close() error {
	// Close services to allow any inflight requests to complete
	// and prevent new requests from being accepted.
	for _, service := range cmd.Services {
		service.Stop()
	}
	defer close(cmd.Closed)
	close(cmd.closing)
	return nil
}

// Options represents the command line options that can be parsed.
type Options struct {
	ConfigPath string
	Action     string
}

var usage = `Runs the cloudive-thumber Listener server.
Usage: cloudive-thumber<cmd> run [flags]
    -config <path>
            Set the path to the configuration file.
            This defaults to the environment variable CLOUDIVE_CONFIG_PATH,
            ~/.cloudive/thumber.conf, or /etc/cloudive/thumber.conf if a file
            is present at any of these locations.
            Disable the automatic loading of a configuration file using
            the null device (such as /dev/null).
`

// ParseFlags parses the command line flags from args and returns an options set.
func (cmd *Command) ParseFlags(args ...string) (Options, error) {
	var options Options
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&options.ConfigPath, "config", "", "")
	fs.StringVar(&options.Action, "s", "", "")
	_ = fs.String("hostname", "", "")
	fs.Usage = func() {
		fmt.Fprintln(cmd.Stderr, usage)
	}
	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}
	return options, nil
}

// ParseConfig parses the config at path.
// Returns a demo configuration if path is blank.
func (cmd *Command) ParseConfig(path string) (*Config, error) {
	// Use demo configuration if no config path is specified.
	if path == "" {
		cmd.Logger.WithField("prefix", "core").Println("no configuration provided, using default settings")
		return NewDemoConfig()
	}

	cmd.Logger.WithField("prefix", "core").Printf("Using configuration at: %s\n", path)

	config := NewConfig()
	if err := config.FromTomlFile(path); err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfigPath returns the config path from the options.
// It will return a path by searching in this order:
//   1. The CLI option in ConfigPath
//   2. The environment variable CLOUDIVE_CONFIG_PATH
//   3. The first thumber.conf file on the path:
//        - ~/.cloudive
//        - /etc/cloudive
func (opt *Options) GetConfigPath() string {
	if opt.ConfigPath != "" {
		if opt.ConfigPath == os.DevNull {
			return ""
		}
		return opt.ConfigPath
	} else if envVar := os.Getenv("CLOUDIVE_CONFIG_PATH"); envVar != "" {
		return envVar
	}

	for _, path := range []string{
		os.ExpandEnv("${HOME}/.cloudive/thumber.conf"),
		"/etc/cloudive/thumber.conf",
	} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
