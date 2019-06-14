package shell

import (
	"fmt"
	"os"

	"mfa4aws/internal/pkg/awssts"
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

// GenerateSTSKeysForExport prints the results for usage on the shell via eval
func GenerateSTSKeysForExport(profile string, tokenCode string) {
	creds, err := awssts.GenerateSTSCredentials(profile, tokenCode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printVars(credentialsToEnvExport(creds))
}

func credentialsToEnvExport(creds *awssts.AWSCredentials) (envVars []string) {
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

func printVars(vars []string) {
	for _, x := range vars {
		fmt.Fprintf(os.Stdout, "%s\n", x)
	}
}
