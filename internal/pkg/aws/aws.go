package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

//Credentials represents the set of attributes used to authenticate to AWS with a short lived session
type Credentials struct {
	AWSAccessKeyID     string        `ini:"aws_access_key_id"`
	AWSSecretAccessKey string        `ini:"aws_secret_access_key"`
	AWSSessionToken    string        `ini:"aws_session_token"`
	AWSSecurityToken   string        `ini:"aws_security_token"`
	PrincipalARN       string        `ini:"x_principal_arn"`
	Expires            time.Duration `ini:"x_security_token_expires"`
}

//GenerateSTSCredentials created STS Credentials
func GenerateSTSCredentials(profile string, tokenCode string) (*Credentials, error) {

	awsSession, err := createSession("", profile)
	if err != nil {
		return nil, err
	}

	iamInstance := iam.New(awsSession)

	mfaSerialNumber, err := getIAMUserMFADevice(iamInstance)
	if err != nil {
		return nil, err
	}

	stsInstance := sts.New(awsSession)
	stsSessionCredentials, err := getSTSSessionToken(stsInstance, tokenCode, mfaSerialNumber)
	if err != nil {
		return nil, err
	}

	identity, err := getSTSIdentity(stsInstance)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		AWSAccessKeyID:     *stsSessionCredentials.AccessKeyId,
		AWSSecretAccessKey: *stsSessionCredentials.SecretAccessKey,
		AWSSessionToken:    *stsSessionCredentials.SessionToken,
		AWSSecurityToken:   *stsSessionCredentials.SessionToken,
		PrincipalARN:       identity.ARN,
		Expires:            time.Until(*stsSessionCredentials.Expiration),
	}, nil
}
