package awssts

import (
	"bytes"
	"fmt"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/spf13/afero"
)

var (
	appFs = afero.NewOsFs()
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

	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(user.HomeDir, ".aws", "credentials")

	f, err := openFile(path)
	if err != nil {
		return nil, ErrAWSCredentialsFileNotFound
	}

	if err := checkProfile(f, profile); err != nil {
		return nil, err
	}

	if err := checkToken(tokenCode); err != nil {
		return nil, err
	}

	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
	}))

	iamInstance := iam.New(awsSession)

	iamUser, err := getIAMUserDetails(iamInstance)
	if err != nil {
		return nil, err
	}

	mfaSerialNumber, err := getIAMUserMFADevice(iamInstance)
	if err != nil {
		return nil, err
	}

	stsInstance := sts.New(awsSession)

	stsSessionCredentials, err := generateSTSSessionCredentials(stsInstance, tokenCode, mfaSerialNumber)
	if err != nil {
		return nil, err
	}

	return &AWSCredentials{
		AWSAccessKeyID: *stsSessionCredentials.AccessKeyId,
		AWSSecretAccessKey: *stsSessionCredentials.SecretAccessKey,
		AWSSessionToken: *stsSessionCredentials.SecretAccessKey,
		AWSSecurityToken: *stsSessionCredentials.SecretAccessKey,
		PrincipalARN: *iamUser.Arn,
		Expires: time.Until(*stsSessionCredentials.Expiration),
	}, nil
}

func checkProfile(file []byte, profile string) error {
	const (
		profileDefault string = "default"

		credentialsAWSAccessKeyID     string = "aws_access_key_id"
		credentialsAWSSecretAccessKey string = "aws_secret_access_key"
	)

	if len(profile) == 0 {
		profile = profileDefault
	}

	creds, err := ini.Load(file)
	if err != nil {
		return ErrInvalidAWSCredentialsFile
	}

	if !creds.Section(profile).HasKey(credentialsAWSAccessKeyID) ||
		!creds.Section(profile).HasKey(credentialsAWSSecretAccessKey) {
		return ErrInvalidAWSCredentialsFile
	}

	return nil
}

func openFile(path string) ([]byte, error) {
	f, err := appFs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func checkToken(token string) error {

	if len(token) <= 5 {
		return ErrInvalidToken
	}

	return nil
}

func getIAMUserDetails(iamInstance iamiface.IAMAPI) (*iam.User, error) {

	identity, err := iamInstance.GetUser(&iam.GetUserInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, fmt.Errorf("Unable to retrive user - %v", aerr.Message())
		}
		return nil, fmt.Errorf("unknown error occurred, %v", err)
	}

	return identity.User, nil
}

func getIAMUserMFADevice(iamInstance iamiface.IAMAPI) (string, error) {
	devices, err := iamInstance.ListMFADevices(&iam.ListMFADevicesInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", fmt.Errorf("Unable to retrive any MFA devices - %v", aerr.Message())
		}
		return "", fmt.Errorf("unknown error occurred, %v", err)
	}

	if len(devices.MFADevices) == 0 {
		return "", ErrNoMFADeviceForUser
	}

	return *devices.MFADevices[0].SerialNumber, nil
}

func generateSTSSessionCredentials(stsInstance stsiface.STSAPI, tokenCode string, mfaDeviceSerialNumber string) (*sts.Credentials, error) {
	stsSession, err := stsInstance.GetSessionToken(&sts.GetSessionTokenInput{
		TokenCode:    &tokenCode,
		SerialNumber: &mfaDeviceSerialNumber,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeExpiredTokenException:
				return nil, ErrTokenHasExpired
			case sts.ErrCodeInvalidIdentityTokenException:
				return nil, ErrInvalidToken
			default:
				return nil, fmt.Errorf("%v For device %s", aerr.Message(), mfaDeviceSerialNumber)
			}
		}
		return nil, fmt.Errorf("unknown error occurred - %v For device %s", err, mfaDeviceSerialNumber)
	}

	return stsSession.Credentials, nil
}
