package smtp

import (
	"fmt"
	"io"
	"mime"
	"net/mail"

	"github.com/nirnanaaa/cloudive-mailer/services/kafka/event"
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

// Deliver performs all necessary operations to send an outgoing email via SMTP
func (s *Service) Deliver(u *event.InboundEmailEvent) error {
	if !s.Config.Enabled {
		return fmt.Errorf("SMTP Service is not enabled, we're not delivering any emails")
	}

	d := NewDialer(s.Config.Hostname, s.Config.Port, s.Config.Username, s.Config.Password)
	m := gomail.NewMessage()
	from := formatEmail(u.Sender.Name, u.Sender.Email)
	to := formatEmail(u.Recipient.Name, u.Recipient.Email)
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", u.Subject)
	htmlData := string(u.Payload)
	m.SetBody("text/html", htmlData)

	for _, attachment := range u.Attachments {
		if s.Config.AttachmentDomainWhitelistEnabled {
			if err := s.CheckAttachmentForDomainWhitelist(attachment.URL); err != nil {
				logrus.WithError(err).Warnf("Attachment Domain whitelisting is enabled and URL does not match whitelist: %s", attachment.URL)
				continue
			}
		}
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

	return d.DialAndSend(m)
}

// SetLogOutput sets the writer to which all logs are written. It must not be
// called after Open is called.
func (s *Service) SetLogOutput(log *logrus.Logger) {
	s.Logger = log
}
