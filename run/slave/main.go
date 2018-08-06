package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nirnanaaa/cloudive-mailer/run/slave/cmd"
	log "github.com/sirupsen/logrus"
)

// These variables are populated via the Go linker.
var (
	version string
	commit  string
	branch  string
)

func init() {
	unknown := "unknown"
	// If commit, branch, or build time are not set, make that clear.
	if version == "" {
		version = unknown
	}
	if commit == "" {
		commit = unknown
	}
	if branch == "" {
		branch = unknown
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	m := NewMain()
	m.Logger = log.New()
	if err := m.Run(os.Args[1:]...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main represents the program execution.
type Main struct {
	Logger *log.Logger

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewMain return a new instance of Main.
func NewMain() *Main {
	return &Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run determines and runs the command specified by the CLI args.
func (m *Main) Run(args ...string) error {
	name, args := ParseCommandName(args)

	// Extract name from args.
	switch name {
	case "", "run":
		cmd := run.NewCommand()

		// Tell the server the build details.
		cmd.Version = version
		cmd.Commit = commit
		cmd.Branch = branch

		if err := cmd.Run(args...); err != nil {
			return fmt.Errorf("run: %s", err)
		}
		cmd.Open()
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
		m.Logger.Println("Listening for signals")
		select {
		case <-signalCh:
			m.Logger.Println("Signal received, initializing clean shutdown...")
			go func() {
				cmd.Close()
			}()
		}

		// Block again until another signal is received, a shutdown timeout elapses,
		// or the Command is gracefully closed
		m.Logger.Println("Waiting for clean shutdown...")
		select {
		case <-signalCh:
			m.Logger.Println("second signal received, initializing hard shutdown")
		case <-time.After(time.Second * 30):
			m.Logger.Println("time limit reached, initializing hard shutdown")
		case <-cmd.Closed:
			m.Logger.Println("server shutdown completed")
		}

	case "config":
		if err := run.NewPrintConfigCommand().Run(args...); err != nil {
			return fmt.Errorf("config: %s", err)
		}
	case "version":
		if err := NewVersionCommand().Run(args...); err != nil {
			return fmt.Errorf("version: %s", err)
		}
	default:
		return fmt.Errorf(`unknown command "%s"`+"\n"+`Run 'roove-thumber help' for usage`+"\n\n", name)
	}

	return nil
}

// ParseCommandName extracts the command name and args from the args list.
func ParseCommandName(args []string) (string, []string) {
	// Retrieve command name as first argument.
	var name string
	if len(args) > 0 {
		if !strings.HasPrefix(args[0], "-") {
			name = args[0]
		} else if args[0] == "-h" || args[0] == "-help" || args[0] == "--help" {
			// Special case -h immediately following binary name
			name = "help"
		}
	}

	// If command is "help" and has an argument then rewrite args to use "-h".
	if name == "help" && len(args) > 2 && !strings.HasPrefix(args[1], "-") {
		return args[1], []string{"-h"}
	}

	// If a named command is specified then return it with its arguments.
	if name != "" {
		return name, args[1:]
	}
	return "", args
}

// VersionCommand represents the command executed by "roove-thumber version".
type VersionCommand struct {
	Stdout io.Writer
	Stderr io.Writer
}

// NewVersionCommand return a new instance of VersionCommand.
func NewVersionCommand() *VersionCommand {
	return &VersionCommand{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run prints the current version and commit info.
func (cmd *VersionCommand) Run(args ...string) error {
	// Parse flags in case -h is specified.
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() { fmt.Fprintln(cmd.Stderr, versionUsage) }
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Print version info.
	fmt.Fprintf(cmd.Stdout, "ROOVE Thumber v%s (git: %s %s)\n", version, branch, commit)

	return nil
}

var versionUsage = `Displays the API version, build branch and git commit hash.
Usage: roove-thumber version
`
