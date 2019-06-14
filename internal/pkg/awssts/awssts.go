package awssts

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

// AWSCredentials represents the set of attributes used to authenticate to AWS with a short lived session
type AWSCredentials struct {
	AWSAccessKeyID     string        `ini:"aws_access_key_id"`
	AWSSecretAccessKey string        `ini:"aws_secret_access_key"`
	AWSSessionToken    string        `ini:"aws_session_token"`
	AWSSecurityToken   string        `ini:"aws_security_token"`
	PrincipalARN       string        `ini:"x_principal_arn"`
	Expires            time.Duration `ini:"x_security_token_expires"`
}

//GenerateSTSCredentials created STS Credentials
func GenerateSTSCredentials(profile string, tokenCode string) (*AWSCredentials, error) {

	result := &AWSCredentials{}

	if err := checkProfile(profile); err != nil {
		return result, err
	}

	if err := checkToken(tokenCode); err != nil {
		return result, err
	}

	sess, err := createNewSession(profile)
	if err != nil {
		return result, err
	}

	iamInstance := iam.New(sess)

	user, err := iamInstance.GetUser(&iam.GetUserInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return result, fmt.Errorf("Unable to retrive user - %v", aerr.Message())
		}
		return result, fmt.Errorf("unknown error occurred, %v", err)
	}

	devices, err := iamInstance.ListMFADevices(&iam.ListMFADevicesInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return result, fmt.Errorf("Unable to retrive any MFA devices - %v", aerr.Message())
		}
		return result, fmt.Errorf("unknown error occurred, %v", err)
	}

	if len(devices.MFADevices) == 0 {
		return result, ErrNoMFADeviceForUser
	}

	sn := devices.MFADevices[0].SerialNumber

	_sts := sts.New(sess)
	stsSession, err := _sts.GetSessionToken(&sts.GetSessionTokenInput{
		TokenCode:    &tokenCode,
		SerialNumber: sn,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeExpiredTokenException:
				return result, ErrTokenHasExpired
			case sts.ErrCodeInvalidIdentityTokenException:
				return result, ErrInvalidToken
			default:
				return result, fmt.Errorf("%v For device %s", aerr.Message(), *sn)
			}
		}
		return result, fmt.Errorf("unknown error occurred - %v For device %s", err, *sn)
	}

	result.AWSAccessKeyID = *stsSession.Credentials.AccessKeyId
	result.AWSSecretAccessKey = *stsSession.Credentials.SecretAccessKey
	result.AWSSessionToken = *stsSession.Credentials.SessionToken
	result.AWSSecurityToken = *stsSession.Credentials.SessionToken
	result.PrincipalARN = *user.User.Arn
	result.Expires = time.Until(*stsSession.Credentials.Expiration)

	return result, nil
}

func checkProfile(profile string) error {
	const (
		profileDefault                string = "default"
		credentialsAWSAccessKeyID     string = "aws_access_key_id"
		credentialsAWSSecretAccessKey string = "aws_secret_access_key"
	)

	if len(profile) == 0 {
		profile = profileDefault
	}

	user, err := user.Current()
	if err != nil {
		return err
	}

	path := filepath.Join(user.HomeDir, ".aws", "credentials")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrAWSCredentialsFileNotFound
	}

	creds, err := ini.Load(path)
	if err != nil {
		return ErrInvalidAWSCredentialsFile
	}

	if !creds.Section(profile).HasKey(credentialsAWSAccessKeyID) ||
		!creds.Section(profile).HasKey(credentialsAWSSecretAccessKey) {
		return ErrInvalidAWSCredentialsFile
	}

	return nil
}

func checkToken(token string) error {

	if len(token) <= 5 {
		return ErrInvalidToken
	}

	return nil
}

func createNewSession(profile string) (*session.Session, error) {

	conf := &aws.Config{
		Credentials: credentials.NewSharedCredentials("", profile),
	}

	sess, err := session.NewSession(conf)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case session.ErrCodeSharedConfig:
				return &session.Session{}, ErrAWSCredentialsFileNotFound
			default:
				return &session.Session{}, fmt.Errorf("Unable to create an AWS session - %v", aerr.Message())
			}
		}
		return &session.Session{}, fmt.Errorf("unknown error occurred, %v", err)
	}
	return sess, nil
}
