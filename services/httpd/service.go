package httpd

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	ln      net.Listener
	addr    string
	err     chan error
	Logger  *logrus.Entry
	Handler *Handler
	Enabled bool
}

// NewService returns a new instance of Service.
func NewService(c Config) *Service {
	s := &Service{
		addr:    c.BindAddress,
		err:     make(chan error),
		Handler: NewHandler(c),
		Enabled: c.Enabled,
		Logger:  logrus.New().WithField("prefix", "httpd"),
	}
	return s
}

// Start starts the service
func (s *Service) Start() error {
	if !s.Enabled {
		s.Logger.Infof("HTTP Service is not enabled. Skipping initialization")
		return nil
	}
	s.Logger.Infof("Starting HTTP service")

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.Logger.Debugln("Listening on HTTP:", listener.Addr().String())
	s.ln = listener

	// wait for the listeners to start
	timeout := time.Now().Add(time.Second)
	for {
		if s.ln.Addr() != nil {
			break
		}

		if time.Now().After(timeout) {
			return fmt.Errorf("unable to open without http listener running")
		}
		time.Sleep(10 * time.Millisecond)
	}
	// Begin listening for requests in a separate goroutine.
	go s.serveTCP()
	return nil
}

// Stop closes the underlying listener.
func (s *Service) Stop() error {
	if s.ln != nil {
		if err := s.ln.Close(); err != nil {
			return err
		}
	}
	return nil
}

// SetLogOutput sets the writer to which all logs are written. It must not be
// called after Open is called.
func (s *Service) SetLogOutput(log *logrus.Logger) {
	l := log.WithField("prefix", "httpd")
	s.Logger = l
	s.Handler.Logger = l
}

// Err returns a channel for fatal errors that occur on the listener.
func (s *Service) Err() <-chan error {
	return s.err
}

// Addr returns the listener's address. Returns nil if listener is closed.
func (s *Service) Addr() net.Addr {
	if s.ln != nil {
		return s.ln.Addr()
	}
	return nil
}

// serveTCP serves the handler from the TCP listener.
func (s *Service) serveTCP() {
	s.serve(s.ln)
}

// serve serves the handler from the listener.
func (s *Service) serve(listener net.Listener) {
	// The listener was closed so exit
	// See https://github.com/golang/go/issues/4373
	err := http.Serve(listener, s.Handler)
	if err != nil && !strings.Contains(err.Error(), "closed") {
		s.err <- fmt.Errorf("listener failed: addr=%s, err=%s", s.Addr(), err)
	}
}
