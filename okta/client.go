package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// Client contains the necessary information needed to communicate with the Okta API.
type Client struct {
	host string
	http *http.Client
}

// NewClient returns a newly configured Client.
func NewClient(host string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		host: host,
		http: &http.Client{
			Jar: jar,
		},
	}, nil
}

// Post makes an HTTP POST request and returns the response as a string.
func (c *Client) Post(path string, payload []byte) (reader io.Reader, err error) {
	url := c.host + path
	resp, err := c.http.Post(url, "application/json", bytes.NewBuffer(payload))
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
	url := c.host + path
	resp, err := c.http.Get(url)
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

	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	errMsg := new(struct {
		Summary string `json:"errorSummary"`
	})

	err := json.NewDecoder(resp.Body).Decode(errMsg)
	if err != nil {
		return err
	}

	return fmt.Errorf("%s (%d)", errMsg.Summary, resp.StatusCode)
}
