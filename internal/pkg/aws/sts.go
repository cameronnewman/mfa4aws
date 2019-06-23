package aws

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

const (
	tokenValidationRegex string = "^[0-9]+$"
)

var (
	tokenValidationRegexComplied = regexp.MustCompilePOSIX(tokenValidationRegex)
)

//STSIdentity represents the STS Identity
type STSIdentity struct {
	Account string
	ARN     string
	UserID  string
}

func validateToken(token string) error {
	if len(token) <= 5 {
		return ErrInvalidToken
	}
	if !tokenValidationRegexComplied.MatchString(token) {
		return ErrInvalidToken
	}

	return nil
}

func getSTSSessionToken(stsInstance stsiface.STSAPI, tokenCode string, mfaDeviceSerialNumber string) (*sts.Credentials, error) {

	if err := validateToken(tokenCode); err != nil {
		return nil, err
	}

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

func getSTSIdentity(stsInstance stsiface.STSAPI) (*STSIdentity, error) {
	identity, err := stsInstance.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, fmt.Errorf("Unable to retrive user - %v", aerr.Message())
		}
		return nil, fmt.Errorf("unknown error occurred, %v", err)
	}

	return &STSIdentity{
		Account: *identity.Account,
		ARN:     *identity.Arn,
		UserID:  *identity.UserId,
	}, nil
}
