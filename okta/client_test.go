package okta

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClient(t *testing.T) {
	host := "https://test.okta.com"
	sid := "test_sid"
	c, err := NewClient(host, sid)
	if err != nil {
		t.Fatalf("unexpected error when creating client: %s", err)
	}

	if c.url.String() != host {
		t.Errorf("got %s, wanted %s", c.url.String(), host)
	}
	if c.http.Jar == nil {
		t.Fatalf("expected HTTP client to have a cookie jar")
	}
	var cookie string
	cookies := c.http.Jar.Cookies(c.url)
	for _, c := range cookies {
		if c.Name == "sid" {
			cookie = c.Value
		}
	}
	if cookie != sid {
		t.Errorf("got %s, wanted %s", cookie, sid)
	}

	if c.http.Jar == nil {
		t.Errorf("expected HTTP client to have a cookie jar")
	}

	msg := "{}"
	errMsg := `{"errorSummary":"error"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(errMsg))
			if err != nil {
				t.Fatalf("unexpected error when writing response: %s", err)
			}
			return
		}
		_, err := w.Write([]byte(msg))
		if err != nil {
			t.Fatalf("unexpected error when writing response: %s", err)
		}
	}))
	defer srv.Close()
	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error when parsing %s: %s", srv.URL, err)
	}
	c.url = u

	reader, err := c.Post("/test", []byte(""))
	if err != nil {
		t.Fatalf("unexpected error when making a successful POST request: %s", err)
	}
	resp, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error when reading a successful POST response: %s", err)
	}
	if string(resp) != msg {
		t.Errorf("got %s, wanted %s", string(resp), msg)
	}

	_, err = c.Post("/error", []byte(""))
	if err == nil || err.Error() != "error (500)" {
		t.Fatalf("expected error when making a failed POST request: %s", err)
	}

	reader, err = c.Get("/test")
	if err != nil {
		t.Fatalf("unexpected error when making a successful GET request: %s", err)
	}
	resp, err = ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error when reading a successful GET response: %s", err)
	}
	if string(resp) != msg {
		t.Errorf("got %s, wanted %s", string(resp), msg)
	}

	_, err = c.Get("/error")
	if err == nil || err.Error() != "error (500)" {
		t.Fatalf("expected error when making a failed GET request: %s", err)
	}
}
