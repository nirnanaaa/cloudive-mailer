package httpd_test

import (
	"testing"

	"github.com/nirnanaaa/cloudive-mailer/services/httpd"
	"github.com/sirupsen/logrus"
)

func CreateService(enabled bool) *httpd.Service {
	c := httpd.NewConfig()
	c.Enabled = enabled
	return httpd.NewService(*c)
}

func TestHTTPD_NewService(t *testing.T) {
	srv := CreateService(true)
	if err := srv.Start(); err != nil {
		t.Fatal(err.Error())
	}
	srv.Stop()
}
func TestHTTPD_NewServiceNotEnabled(t *testing.T) {
	srv := CreateService(false)
	if err := srv.Start(); err != nil {
		t.Fatal(err.Error())
	}
	srv.Stop()
}
func TestHTTPD_SetLogOutput(t *testing.T) {
	srv := CreateService(true)
	srv.SetLogOutput(logrus.New())
	if err := srv.Start(); err != nil {
		t.Fatal(err.Error())
	}
	srv.Stop()
}
