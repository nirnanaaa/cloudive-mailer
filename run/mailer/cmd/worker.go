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
	"github.com/nirnanaaa/cloudive-mailer/services/smtp"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Service is an interface definition for each service.
type WorkerService interface {
	Start() error
	Stop() error
}

// Command represents the command executed by "roove-thumber run".
type WorkerCommand struct {
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
func WorkerNewCommand() *WorkerCommand {
	logger := log.New()
	logger.Formatter = new(prefixed.TextFormatter)
	return &WorkerCommand{
		closing: make(chan struct{}),
		err:     make(chan error),
		Closed:  make(chan struct{}),
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Logger:  logger,
	}
}

// Run parses the config from args and runs the server.
func (cmd *WorkerCommand) WorkerRun(args ...string) error {
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
		return fmt.Errorf("%s. To generate a valid configuration file run `multimedia-worker config > multimedia-worker.generated.conf`", err)
	}
	level, err := log.ParseLevel(config.Meta.LogLevel)
	if err != nil {
		return err
	}
	logger.SetLevel(level)
	flag.Parse()
	smtpService := smtp.NewService(config.SMTP)
	kafkaService := kafka.NewService(config.Kafka)
	kafkaService.SetLogOutput(logger, "[kafka]")
	kafkaService.SetDefaultMessageProcessor(smtpService)
	httpdService := httpd.NewService(*config.HTTPD)
	httpdService.SetLogOutput(logger)
	smtpService.SetLogOutput(logger)

	cmd.Services = append(cmd.Services, httpdService)
	cmd.Services = append(cmd.Services, kafkaService)
	cmd.Services = append(cmd.Services, smtpService)
	return nil
}

// Open opens all services
func (cmd *WorkerCommand) WorkerOpen() error {
	for _, service := range cmd.Services {
		if err := service.Start(); err != nil {
			log.Errorf(err.Error())
			return fmt.Errorf("open service: %s", err)
		}
	}
	return nil
}

// Close shuts down the server.
func (cmd *WorkerCommand) WorkerClose() error {
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
type WorkerOptions struct {
	ConfigPath string
	Action     string
}

var workerUsage = `Runs the roove-thumber Listener server.
Usage: roove-thumber<cmd> run [flags]
    -config <path>
            Set the path to the configuration file.
            This defaults to the environment variable ROOVE_CONFIG_PATH,
            ~/.roove/thumber.conf, or /etc/roove/thumber.conf if a file
            is present at any of these locations.
            Disable the automatic loading of a configuration file using
            the null device (such as /dev/null).
`

// ParseFlags parses the command line flags from args and returns an options set.
func (cmd *WorkerCommand) ParseFlags(args ...string) (Options, error) {
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
func (cmd *WorkerCommand) ParseConfig(path string) (*Config, error) {
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
