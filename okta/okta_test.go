package okta

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/test"
)

func TestLogin(t *testing.T) {
	appPath := "/app/url"
	appSuccessResponse := test.LoadTestFile(t, "app_success_response.html")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == appPath {
			_, err := w.Write([]byte(appSuccessResponse))
			if err != nil {
				t.Fatalf("unexpected error when writing response: %s", err)
			}
			return
		}
		_, err := w.Write([]byte("{}"))
		if err != nil {
			t.Fatalf("unexpected error when writing response: %s", err)
		}
	}))
	defer srv.Close()
	conf, err := config.New("")
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}
	conf.OktaHost = srv.URL
	conf.OktaAppPath = appPath

	i := test.NewNoopInput()
	_, _, err = Login(conf, i, "", "")
	if err != nil {
		t.Fatalf("unexpected error when logging in: %s", err)
	}
}
