package smtp

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/mail"
	"time"

	"github.com/sirupsen/logrus"

	gomail "gopkg.in/gomail.v2"
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	Logger *logrus.Logger
	Config Config
}

// NewService returns a new instance of Service.
func NewService(c Config) *Service {
	s := &Service{
		Config: c,
	}
	return s
}

// Start starts the service
func (s *Service) Start() error {
	return nil
}

// Stop closes the underlying listener.
func (s *Service) Stop() error {
	return nil
}

func formatEmail(name, email string) string {
	s := fmt.Sprintf("%s <%s>", name, email)
	se := mime.QEncoding.Encode("utf-8", s)
	_, e := mail.ParseAddress(se)
	if e != nil {
		return email
	}
	return s
}

func (s *Service) DownloadAttachment(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
func (s *Service) Deliver(u *OutboundEmailEvent, tries int) {
	if !s.Config.Enabled {
		return
	}
	if tries > 10 {
		return
	}

	d := NewDialer(s.Config.Hostname, s.Config.Port, s.Config.Username, s.Config.Password)
	m := gomail.NewMessage()
	from := formatEmail(s.Config.FromName, s.Config.FromMail)
	to := formatEmail(u.Recipient.Name, u.Recipient.Email)
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", u.Subject)
	htmlData := string(u.Payload)
	m.SetBody("text/html", htmlData)

	for _, attachment := range u.Attachments {
		m.Attach(attachment.Name, gomail.SetCopyFunc(func(w io.Writer) error {
			r, err := s.DownloadAttachment(attachment.URL)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, r); err != nil {
				return err
			}
			return nil
		}))
	}

	if err := d.DialAndSend(m); err != nil {
		time.Sleep(time.Duration(tries*tries) * time.Second)
		logrus.WithError(err).Errorf("trying next delivery in %ds", tries*tries)
		s.Deliver(u, tries+1)
		return
	}
}

// SetLogOutput sets the writer to which all logs are written. It must not be
// called after Open is called.
func (s *Service) SetLogOutput(log *logrus.Logger) {
	s.Logger = log
}
