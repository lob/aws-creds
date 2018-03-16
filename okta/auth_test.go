package okta

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/test"
)

func TestAuthLogin(t *testing.T) {
	path := path.Join(os.TempDir(), "aws-creds", "config")
	defer cleanup(t, path)
	conf := config.New(path)

	loginSuccessResponse := loadTestFile(t, "login_success_response.json")
	verifySuccessResponse := loadTestFile(t, "verify_success_response.json")
	handler := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/api/v1/authn" {
			_, err := w.Write([]byte(loginSuccessResponse))
			if err != nil {
				t.Fatalf("unexpected error when writing response: %s", err)
			}
			return
		}
		if strings.Contains(req.URL.Path, "/verify") {
			_, err := w.Write([]byte(verifySuccessResponse))
			if err != nil {
				t.Fatalf("unexpected error when writing response: %s", err)
			}
			return
		}
	}
	srv, c := testServerAndClient(t, handler)
	defer srv.Close()

	auth, err := login(c, "", "")
	if err != nil {
		t.Fatalf("unexpected error when logging in: %s", err)
	}

	if auth.Status != "MFA_REQUIRED" {
		t.Errorf("got %s, wanted %s", auth.Status, "MFA_REQUIRED")
	}
	if len(auth.Embedded.Factors) != 2 {
		t.Errorf("got len(auth.Embedded.Factors) = %d, wanted %d", len(auth.Embedded.Factors), 2)
	}

	cases := []struct {
		reason, factorType string
		responses          []string
		err                bool
	}{
		{"with no preferred factor", "", []string{"0", "n", "123456", ""}, false},
		{"with an out-of-bound factor index", "", []string{"3", "0", "n", "123456", ""}, false},
		{"with an invalid factor index", "", []string{"err", "0", "n", "123456", ""}, false},
		{"with an unsupported factor", "", []string{"1", "n", ""}, true},
		{"and saving the preferred factor", "", []string{"0", "y", "123456", ""}, false},
		{"with a preferred factor", "token:software:totp", []string{"123456", ""}, false},
		{"with an invalid preferred factor", "invalid", []string{""}, true},
	}

	for _, tc := range cases {
		conf.PreferredFactorType = tc.factorType
		i := test.NewArrayInput(tc.responses)
		err = auth.verifyMFA(c, conf, i)
		if tc.err && err == nil {
			t.Errorf("expected error when verifying MFA %s", tc.reason)
		} else if !tc.err && err != nil {
			t.Errorf("unexpected error when verifying MFA %s: %s", tc.reason, err)
		}
	}
}

func loadTestFile(t *testing.T, name string) string {
	path := fmt.Sprintf("testdata/%s", name)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error when reading file %s: %s", path, err)
	}
	return string(contents)
}

func testServerAndClient(t *testing.T, handler func(http.ResponseWriter, *http.Request)) (*httptest.Server, *Client) {
	srv := httptest.NewServer(http.HandlerFunc(handler))
	c, err := NewClient("test")
	if err != nil {
		t.Fatalf("unexpected error when creating client: %s", err)
	}
	c.host = srv.URL
	return srv, c
}

func cleanup(t *testing.T, path string) {
	dir := filepath.Dir(path)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("unexpected error when cleaning up: %s", err)
	}
}
