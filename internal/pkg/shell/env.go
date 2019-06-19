package shell

import (
	"fmt"
	"mfa4aws/internal/pkg/aws"
)

const (
	envNameAWSAccessKey     string = "AWS_ACCESS_KEY_ID"
	envNameAWSSecretKey     string = "AWS_SECRET_ACCESS_KEY"
	envNameAWSSessionToken  string = "AWS_SESSION_TOKEN"
	envNameAWSSecurityToken string = "AWS_SECURITY_TOKEN"
	envNameXPrincipalARN    string = "X_PRINCIPAL_ARN"
	envNameExpires          string = "EXPIRES"

	bashExport string = "export"
)

//BuildEnvVars - constructs a string array from the Credentials
func BuildEnvVars(creds *aws.Credentials) (envVars []string) {
	envVars = append(envVars,
		fmt.Sprintf("%s %s=%s", bashExport, envNameAWSAccessKey, creds.AWSAccessKeyID),
		fmt.Sprintf("%s %s=%s", bashExport, envNameAWSSecretKey, creds.AWSSecretAccessKey),
		fmt.Sprintf("%s %s=%s", bashExport, envNameAWSSessionToken, creds.AWSSessionToken),
		fmt.Sprintf("%s %s=%s", bashExport, envNameAWSSecurityToken, creds.AWSSecurityToken),
		fmt.Sprintf("%s %s=%s", bashExport, envNameXPrincipalARN, creds.PrincipalARN),
		fmt.Sprintf("%s %s=%s", bashExport, envNameExpires, creds.Expires.String()),
	)

	return envVars
}
