package aws

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/lob/aws-creds/internal/test"
	"github.com/lob/aws-creds/pkg/config"
	"github.com/lob/aws-creds/pkg/okta"
)

func TestGetCreds(t *testing.T) {
	roleARN := "role_arn"
	duration := "1800"
	creds := test.NewCredentials()
	svc := &test.MockSTS{Creds: creds}
	saml := &okta.SAMLResponse{
		Attributes: []okta.Attribute{
			{
				Name:   roleSAMLAttribute,
				Values: []string{fmt.Sprintf("principal_arn,%s", roleARN)},
			},
			{
				Name:   durationSAMLAttribute,
				Values: []string{duration},
			},
		},
	}
	profile := &config.Profile{Name: "staging", RoleARN: roleARN}

	c, err := GetCreds(svc, saml, profile)
	if err != nil {
		t.Fatalf("unexpected error when getting creds: %s", err)
	}

	cases := []struct {
		got, want string
	}{
		{*c.AccessKeyId, *creds.AccessKeyId},
		{*c.SecretAccessKey, *creds.SecretAccessKey},
		{*c.SessionToken, *creds.SessionToken},
		{strconv.Itoa(int(svc.Duration)), duration},
	}

	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("got %s, wanted %s", tc.got, tc.want)
		}
	}

	badProfile := &config.Profile{Name: "staging", RoleARN: "invalid"}
	_, err = GetCreds(svc, saml, badProfile)
	if err == nil {
		t.Errorf("expected error when getting creds for an invalid profile")
	}
}

func TestGetCredsNilRoleARN(t *testing.T) {
	roleARN := ""
	duration := "1800"
	creds := test.NewCredentials()
	svc := &test.MockSTS{Creds: creds}
	saml := &okta.SAMLResponse{
		Attributes: []okta.Attribute{
			{
				Name:   roleSAMLAttribute,
				Values: []string{fmt.Sprintf("principal_arn,%s", roleARN)},
			},
			{
				Name:   durationSAMLAttribute,
				Values: []string{duration},
			},
		},
	}
	profile := &config.Profile{Name: "staging", RoleARN: roleARN}

	_, err := GetCreds(svc, saml, profile)
	if err == nil {
		t.Fatalf("expected error when profile ARN is empty string")
	}
}
