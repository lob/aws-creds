package okta

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lob/aws-creds/config"
)

func TestLogin(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("{}"))
		if err != nil {
			t.Fatalf("unexpected error when writing response: %s", err)
		}
	}))
	defer srv.Close()
	conf := config.New("")
	conf.OktaHost = srv.URL

	err := Login(conf, "")
	if err != nil {
		t.Fatalf("unexpected error when logging in: %s", err)
	}
}
