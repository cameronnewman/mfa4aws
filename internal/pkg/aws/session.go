package aws

import (
	"bytes"
	"os/user"
	"path/filepath"

	"gopkg.in/ini.v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/afero"
)

var (
	appFs = afero.NewOsFs()
)

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

func createSession(profile string) (*session.Session, error) {
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

	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewSharedCredentials(path, profile),
		},
	}))

	return awsSession, nil
}
