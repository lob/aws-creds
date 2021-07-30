package aws

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/lob/aws-creds/pkg/config"
	"github.com/lob/aws-creds/pkg/okta"
)

var (
	roleSAMLAttribute           = "https://aws.amazon.com/SAML/Attributes/Role"
	durationSAMLAttribute       = "https://aws.amazon.com/SAML/Attributes/SessionDuration"
	defaultDuration       int64 = 3600
)

// GetCreds fetches AWS credentials using a SAML response.
func GetCreds(svc stsiface.STSAPI, saml *okta.SAMLResponse, profile *config.Profile) (*sts.Credentials, error) {
	roles, duration := parseSAMLAttributes(saml)

	// Override default duration
	if profile.Duration > 3600 {
		duration = profile.Duration
	}

	var role string
	for _, r := range roles {
		if strings.Contains(r, profile.RoleARN) {
			role = r
			break
		}
	}
	if role == "" {
		return nil, fmt.Errorf("%s is not a valid role you can assume", profile.RoleARN)
	}

	arns := strings.Split(role, ",")
	principalARN := arns[0]
	roleARN := arns[1]

	params := &sts.AssumeRoleWithSAMLInput{
		PrincipalArn:    aws.String(principalARN),
		RoleArn:         aws.String(roleARN),
		SAMLAssertion:   aws.String(saml.Raw),
		DurationSeconds: aws.Int64(duration),
	}

	resp, err := svc.AssumeRoleWithSAML(params)
	if err != nil {
		fmt.Println("Failed to assume role for " + profile.Name)
		return nil, err
	}

	return resp.Credentials, nil
}

func parseSAMLAttributes(saml *okta.SAMLResponse) ([]string, int64) {
	var roles []string
	duration := defaultDuration

	for _, attr := range saml.Attributes {
		switch attr.Name {
		case roleSAMLAttribute:
			roles = attr.Values
		case durationSAMLAttribute:
			if len(attr.Values) > 0 {
				d, err := strconv.Atoi(attr.Values[0])
				if err == nil && duration > int64(d) {
					duration = int64(d)
				}
			}
		}
	}

	return roles, duration
}
