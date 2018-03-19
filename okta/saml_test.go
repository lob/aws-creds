package okta

import (
	"net/http"
	"testing"

	"github.com/lob/aws-creds/test"
)

func TestGetSAMLResponse(t *testing.T) {
	token := "session_token"
	appSuccessResponse := test.LoadTestFile(t, "app_success_response.html")
	appFailureResponse := test.LoadTestFile(t, "app_failure_response.html")
	handler := func(w http.ResponseWriter, r *http.Request) {
		tkn := r.URL.Query()["onetimetoken"][0]
		if tkn == token {
			_, err := w.Write([]byte(appSuccessResponse))
			if err != nil {
				t.Fatalf("unexpected error when writing response: %s", err)
			}
			return
		}
		_, err := w.Write([]byte(appFailureResponse))
		if err != nil {
			t.Fatalf("unexpected error when writing response: %s", err)
		}
	}
	srv, c := testServerAndClient(t, handler)
	defer srv.Close()

	resp, err := getSAMLResponse(c, "", token)
	if err != nil {
		t.Fatalf("unexpected error when getting SAML response: %s", err)
	}
	if resp.Raw == "" {
		t.Errorf("expected raw SAML response to not be empty")
	}

	_, err = getSAMLResponse(c, "", "invalid")
	if err == nil {
		t.Fatalf("expected error when getting SAML response with an invalid token")
	}
}
