package smtp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// DownloadAttachment simply downloads an url into a io.Reader format
func (s *Service) DownloadAttachment(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// CheckAttachmentForDomainWhitelist checks if a requesting domain is whitelisted inside our config.
func (s *Service) CheckAttachmentForDomainWhitelist(inputURL string) error {
	u, err := url.Parse(inputURL)
	if err != nil {
		return err
	}
	hostname := u.Hostname()
	for _, entry := range s.Config.DomainWhitelist {
		if hostname == entry {
			return nil
		}
	}
	return fmt.Errorf("Domain %s is not whitelisted", hostname)
}
