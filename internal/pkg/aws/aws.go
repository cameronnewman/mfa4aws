package aws

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

	"github.com/spf13/afero"
)

var (
	appFs = afero.NewOsFs()
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

	const (
		awsCredentialsFolder string = ".aws"
		awsCredentialsFile   string = "credentials"
	)

	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(user.HomeDir, awsCredentialsFolder, awsCredentialsFile)

	f, err := openFile(path)
	if err != nil {
		return nil, ErrAWSCredentialsFileNotFound
	}

	if err := validateProfile(f, profile); err != nil {
		return nil, err
	}

	if err := validateToken(tokenCode); err != nil {
		return nil, err
	}

	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
	}))

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
		AWSSessionToken:    *stsSessionCredentials.SecretAccessKey,
		AWSSecurityToken:   *stsSessionCredentials.SecretAccessKey,
		PrincipalARN:       identity.ARN,
		Expires:            time.Until(*stsSessionCredentials.Expiration),
	}, nil
}

func validateProfile(file []byte, profile string) error {
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

func validateToken(token string) error {
	if len(token) <= 5 {
		return ErrInvalidToken
	}

	return nil
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