package aws

import "errors"

var (
	//ErrAWSCredentialsFileNotFound return when no AWS credentials file can be found at $HOME/.aws/credentials
	ErrAWSCredentialsFileNotFound = errors.New("AWS Credentials file not found at $HOME/.aws/credentials")

	//ErrInvalidAWSCredentialsFile return when AWS credentials file is invaild
	ErrInvalidAWSCredentialsFile = errors.New("AWS Credentials at $HOME/.aws/credentials is invalid")

	//ErrNoMFADeviceForUser is return when no MFA devices have been found for the user
	ErrNoMFADeviceForUser = errors.New("No MFA devices configured for user")

	//ErrTokenHasExpired is returned when the given token has expired
	ErrTokenHasExpired = errors.New("Token has expired")

	//ErrInvalidToken is returned when an invalid token is supplied
	ErrInvalidToken = errors.New("Invalid token code")
)
