package okta

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthLogin(t *testing.T) {
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
