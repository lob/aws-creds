package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/test"
	"github.com/zalando/go-keyring"
)

func TestExecuteRefresh(t *testing.T) {
	keyring.MockInit()
	password := "password"
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
	cfp := path.Join(os.TempDir(), ".aws", "credentials")
	test.PrepTempFile(t, cfp)
	conf, err := config.New("")
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}
	conf.Username = "user"
	conf.OktaHost = srv.URL
	conf.OktaAppPath = appPath
	conf.Profiles = []*config.Profile{{"staging", "arn:aws:iam::123456789001:role/EngineeringRole"}}
	conf.CredentialsFilepath = cfp
	defer test.Cleanup(t, conf.CredentialsFilepath)
	creds := test.NewCredentials()
	cmd := &Cmd{
		Command: "",
		Config:  conf,
		Profile: conf.Profiles[0].Name,
		Input:   test.NewNoopInput(),
		STS:     &test.MockSTS{Creds: creds},
	}

	err = executeRefresh(cmd)
	if err != nil {
		t.Fatalf("unexpected error when executing refresh: %s", err)
	}

	cmd.Profile = "invalid"
	err = executeRefresh(cmd)
	if err == nil {
		t.Fatalf("expected error when executing refresh with an invalid profile: %s", err)
	}

	cmd.Profile = ""
	err = executeRefresh(cmd)
	if err == nil {
		t.Fatalf("expected error when executing refresh without a profile: %s", err)
	}

	cmd.Profile = conf.Profiles[0].Name
	err = keyring.Set(keyringService, conf.Username, password)
	if err != nil {
		t.Fatalf("unexpected error when setting password in keyring: %s", err)
	}
	err = executeRefresh(cmd)
	if err != nil {
		t.Errorf("unexpected error when executing refresh with a saved password: %s", err)
	}
	err = keyring.Delete(keyringService, conf.Username)
	if err != nil {
		t.Fatalf("unexpected error when deleting password in keyring: %s", err)
	}

	cmd.Input = test.NewArrayInput([]string{password, "y"})
	err = executeRefresh(cmd)
	if err != nil {
		t.Errorf("unexpected error when executing refresh when trying to save password: %s", err)
	}
	p, err := keyring.Get(keyringService, conf.Username)
	if err != nil {
		t.Errorf("unexpected error when getting password from keyring: %s", err)
	}
	if p != password {
		t.Errorf("got %s, wanted %s", p, password)
	}
}
