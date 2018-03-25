package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// Client contains the necessary information needed to communicate with the Okta API.
type Client struct {
	url  *url.URL
	http *http.Client
}

// NewClient returns a newly configured Client with the given host and possible Okta
// session cookie. If the session cookie is empty, there was no session previously saved.
func NewClient(host, sessionCookie string) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	if sessionCookie != "" {
		jar.SetCookies(u, []*http.Cookie{
			{
				Name:  "sid",
				Value: sessionCookie,
			},
		})
	}

	return &Client{
		url: u,
		http: &http.Client{
			Jar: jar,
		},
	}, nil
}

// Post makes an HTTP POST request and returns the response as a string.
func (c *Client) Post(path string, payload []byte) (reader io.Reader, err error) {
	u := c.url.String() + path
	resp, err := c.http.Post(u, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	defer func() {
		e := resp.Body.Close()
		if err == nil {
			err = e
		}
	}()

	err = checkError(resp)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	reader = strings.NewReader(string(body))
	return
}

// Get makes an HTTP GET request and returns the response as a string.
func (c *Client) Get(path string) (reader io.Reader, err error) {
	u := c.url.String() + path
	resp, err := c.http.Get(u)
	if err != nil {
		return
	}

	defer func() {
		e := resp.Body.Close()
		if err == nil {
			err = e
		}
	}()

	err = checkError(resp)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	reader = strings.NewReader(string(body))
	return
}

func checkError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	errMsg := new(struct {
		Summary string `json:"errorSummary"`
	})

	err := json.NewDecoder(resp.Body).Decode(errMsg)
	if err != nil {
		return err
	}

	return fmt.Errorf("%s (%d)", errMsg.Summary, resp.StatusCode)
}
