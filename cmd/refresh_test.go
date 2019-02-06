package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/test"
	"github.com/zalando/go-keyring"
)

type authPayload struct {
	Username string `json:"username"`
}

func TestExecuteRefresh(t *testing.T) {
	keyring.MockInit()
	mfaUser := "mfa"
	password := "password"
	appPath := "/app/url"
	appSuccessResponse := test.LoadTestFile(t, "app_success_response.html")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		resp := []byte("{}")
		if req.URL.Path == appPath {
			resp = []byte(appSuccessResponse)
		}
		if req.URL.Path == "/api/v1/authn" {
			body := &authPayload{}
			err := json.NewDecoder(req.Body).Decode(body)
			if err != nil {
				t.Fatalf("unexpected error when decoding body: %s", err)
			}
			if body.Username == mfaUser {
				resp = []byte(`{"status":"MFA_REQUIRED"}`)
			}
		}
		_, err := w.Write(resp)
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
		Command:  "",
		Config:   conf,
		Profiles: []string{conf.Profiles[0].Name},
		Input:    test.NewNoopInput(),
		STS:      &test.MockSTS{Creds: creds},
	}

	err = executeRefresh(cmd)
	if err != nil {
		t.Fatalf("unexpected error when executing refresh: %s", err)
	}

	cmd.Profiles = []string{"invalid"}
	err = executeRefresh(cmd)
	if err == nil {
		t.Fatalf("expected error when executing refresh with an invalid profile: %s", err)
	}

	cmd.Profiles = []string{}
	err = executeRefresh(cmd)
	if err == nil {
		t.Fatalf("expected error when executing refresh without a profile: %s", err)
	}

	cmd.Profiles = []string{conf.Profiles[0].Name}
	err = keyring.Set(keyringPasswordService, conf.Username, password)
	if err != nil {
		t.Fatalf("unexpected error when setting password in keyring: %s", err)
	}
	err = executeRefresh(cmd)
	if err != nil {
		t.Errorf("unexpected error when executing refresh with a saved password: %s", err)
	}
	conf.Username = mfaUser
	conf.PreferredFactorType = "invalid"
	err = keyring.Set(keyringPasswordService, conf.Username, password)
	if err != nil {
		t.Fatalf("unexpected error when setting password in keyring: %s", err)
	}
	err = executeRefresh(cmd)
	if err != nil && !strings.Contains(err.Error(), "MFA") {
		t.Errorf("unexpected error when executing refresh to delete saved password on err: %s", err)
	}
	_, err = keyring.Get(keyringPasswordService, conf.Username)
	if err != keyring.ErrNotFound {
		t.Fatalf("expected not found error when getting deleted password: %s", err)
	}

	conf.Username = "user"
	conf.PreferredFactorType = ""
	cmd.Input = test.NewArrayInput([]string{password, "y"})
	err = executeRefresh(cmd)
	if err != nil {
		t.Errorf("unexpected error when executing refresh when trying to save password: %s", err)
	}
	p, err := keyring.Get(keyringPasswordService, conf.Username)
	if err != nil {
		t.Errorf("unexpected error when getting password from keyring: %s", err)
	}
	if p != password {
		t.Errorf("got %s, wanted %s", p, password)
	}
}
