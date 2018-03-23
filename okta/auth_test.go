package okta

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
)

type verifyPayload struct {
	StateToken string `json:"stateToken"`
	Answer     string `json:"answer"`
}

const (
	successCode = "123456"
	errorCode   = "000000"
	emptyCpde   = ""
)

func TestAuthLogin(t *testing.T) {
	path := path.Join(os.TempDir(), "aws-creds", "config")
	defer test.Cleanup(t, path)
	conf, err := config.New(path)
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}

	srv, c := testServerAndClient(t, authServerHandler(t))
	defer srv.Close()

	auth, err := login(c, "", "")
	if err != nil {
		t.Fatalf("unexpected error when logging in: %s", err)
	}

	if auth.Status != "MFA_REQUIRED" {
		t.Errorf("got %s, wanted %s", auth.Status, "MFA_REQUIRED")
	}
	if len(auth.Embedded.Factors) != 3 {
		t.Errorf("got len(auth.Embedded.Factors) = %d, wanted %d", len(auth.Embedded.Factors), 3)
	}

	cases := []struct {
		reason, factorType string
		responses          []string
		err                bool
	}{
		{"with no preferred factor", "", []string{"0", "n", successCode, ""}, false},
		{"with an out-of-bound factor index", "", []string{"3", "0", "n", successCode, ""}, false},
		{"with an invalid factor index", "", []string{"err", "0", "n", successCode, ""}, false},
		{"with an unsupported factor", "", []string{"2", "n", ""}, true},
		{"with a failed attempt", "", []string{"0", "n", errorCode, successCode, ""}, false},
		{"and saving the preferred factor", "", []string{"0", "y", successCode, ""}, false},
		{"with a preferred factor", "token:software:totp", []string{successCode, ""}, false},
		{"with SMS", "", []string{"1", "n", successCode, ""}, false},
		{"with SMS and a failed attempt", "", []string{"1", "n", errorCode, successCode, ""}, false},
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

func authServerHandler(t *testing.T) http.HandlerFunc {
	loginSuccessResponse := test.LoadTestFile(t, "login_success_response.json")
	verifySuccessResponse := test.LoadTestFile(t, "verify_success_response.json")
	verifyErrorResponse := test.LoadTestFile(t, "verify_error_response.json")
	return func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/api/v1/authn" {
			_, err := w.Write([]byte(loginSuccessResponse))
			if err != nil {
				t.Fatalf("unexpected error when writing response: %s", err)
			}
			return
		}
		if strings.Contains(req.URL.Path, "/verify") {
			body := &verifyPayload{}
			err := json.NewDecoder(req.Body).Decode(body)
			if err != nil {
				t.Fatalf("unexpected error when decoding request body: %s", err)
			}
			var resp string
			switch body.Answer {
			case successCode:
				resp = verifySuccessResponse
			case emptyCpde:
				resp = verifySuccessResponse
			case errorCode:
				w.WriteHeader(403)
				resp = verifyErrorResponse
			}
			_, err = w.Write([]byte(resp))
			if err != nil {
				t.Fatalf("unexpected error when writing response: %s", err)
			}
			return
		}
	}
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
