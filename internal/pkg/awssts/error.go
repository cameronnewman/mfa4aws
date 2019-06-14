package awssts

import "golang.org/x/xerrors"

var (
	//ErrAWSCredentialsFileNotFound return when no AWS credentials file can be found at $HOME/.aws/credentials
	ErrAWSCredentialsFileNotFound = xerrors.New("AWS Credentials file not found at $HOME/.aws/credentials")

	//ErrInvalidAWSCredentialsFile return when AWS credentials file is invaild
	ErrInvalidAWSCredentialsFile = xerrors.New("AWS Credentials at $HOME/.aws/credentials is invalid")

	//ErrNoMFADeviceForUser is return when no MFA devices have been found for the user
	ErrNoMFADeviceForUser = xerrors.New("No MFA devices configured for user")

	//ErrTokenHasExpired is returned when the given token has expired
	ErrTokenHasExpired = xerrors.New("Token has expired")

	//ErrInvalidToken is returned when an invalid token is supplied
	ErrInvalidToken = xerrors.New("Invalid token code")
)
