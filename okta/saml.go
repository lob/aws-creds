package okta

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"

	"golang.org/x/net/html"
)

// SAMLResponse contains the parsed SAML assertion, along with the raw base64 encoded response.
type SAMLResponse struct {
	XMLName    xml.Name    `xml:"Response"`
	Attributes []Attribute `xml:"Assertion>AttributeStatement>Attribute"`
	Raw        string
}

// Attribute contains SAML attributes.
type Attribute struct {
	Name   string   `xml:"Name,attr"`
	Values []string `xml:"AttributeValue"`
}

func getSAMLResponse(c *Client, appPath, sessionToken string) (*SAMLResponse, error) {
	url := fmt.Sprintf("%s?onetimetoken=%s", appPath, sessionToken)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(resp)
	if err != nil {
		return nil, err
	}
	input := getInputNode(doc)
	if input == nil {
		return nil, errors.New("SAML assertion not found")
	}
	encoded := getValue(input)

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	saml := &SAMLResponse{Raw: encoded}
	err = xml.Unmarshal(decoded, saml)
	if err != nil {
		return nil, err
	}
	return saml, nil
}

// getInputNode initially takes in the <html> node of the page and traverses the DOM
// to find <input name="SAMLResponse"> which contains the base64 encoded SAML response.
func getInputNode(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "input" {
		for _, attr := range node.Attr {
			if attr.Key == "name" && attr.Val == "SAMLResponse" {
				return node
			}
		}
	}

	var input *html.Node
	for c := node.FirstChild; input == nil && c != nil; c = c.NextSibling {
		input = getInputNode(c)
	}
	return input
}

// getValue retrieves the value attribute of the given node.
func getValue(input *html.Node) string {
	var val string
	for _, attr := range input.Attr {
		if attr.Key == "value" {
			val = attr.Val
		}
	}
	return val
}
