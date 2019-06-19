package aws

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"mfa4aws/internal/pkg/aws/mock/stsmock"

	"golang.org/x/xerrors"
)

func Test_getSTSSessionToken(t *testing.T) {
	type args struct {
		stsInstance           stsiface.STSAPI
		tokenCode             string
		mfaDeviceSerialNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    *sts.Credentials
		wantErr bool
	}{
		{
			"Vaild/EmptyResult",
			args{
				stsInstance: &stsmock.STSAPIMock{
					GetSessionTokenFunc: func(in1 *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
						return &sts.GetSessionTokenOutput{}, nil
					},
				},
				tokenCode:             "123456",
				mfaDeviceSerialNumber: "sfagstfey",
			},
			nil,
			false,
		},
		{
			"Invaild/awserrError",
			args{
				stsInstance: &stsmock.STSAPIMock{
					GetSessionTokenFunc: func(in1 *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
						return nil, awserr.New("5000", "blah", xerrors.New("blah"))
					},
				},
				tokenCode:             "123456",
				mfaDeviceSerialNumber: "sfagstfey",
			},
			nil,
			true,
		},
		{
			"Invaild/awserrError/ErrCodeExpiredTokenException",
			args{
				stsInstance: &stsmock.STSAPIMock{
					GetSessionTokenFunc: func(in1 *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
						return nil, awserr.New(sts.ErrCodeExpiredTokenException, "Blah", xerrors.New("blah"))
					},
				},
				tokenCode:             "123456",
				mfaDeviceSerialNumber: "sfagstfey",
			},
			nil,
			true,
		},
		{
			"Invaild/awserrError/ErrCodeInvalidIdentityTokenException",
			args{
				stsInstance: &stsmock.STSAPIMock{
					GetSessionTokenFunc: func(in1 *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
						return nil, awserr.New(sts.ErrCodeInvalidIdentityTokenException, "Blah", xerrors.New("blah"))
					},
				},
				tokenCode:             "123456",
				mfaDeviceSerialNumber: "sfagstfey",
			},
			nil,
			true,
		},
		{
			"Invaild/Error",
			args{
				stsInstance: &stsmock.STSAPIMock{
					GetSessionTokenFunc: func(in1 *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
						return nil, xerrors.New("blah")
					},
				},
				tokenCode:             "123456",
				mfaDeviceSerialNumber: "sfagstfey",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSTSSessionToken(tt.args.stsInstance, tt.args.tokenCode, tt.args.mfaDeviceSerialNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSTSSessionToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSTSSessionToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSTSIdentity(t *testing.T) {
	type args struct {
		stsInstance stsiface.STSAPI
	}
	tests := []struct {
		name    string
		args    args
		want    *STSIdentity
		wantErr bool
	}{
		{
			"Valid/User",
			args{
				stsInstance: &stsmock.STSAPIMock{
					GetCallerIdentityFunc: func(in1 *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {

						account := "342563637373"
						arn := "ashgajsdhgajsdg"
						userID := "asjkdhkasdhaksd"

						return &sts.GetCallerIdentityOutput{
							Account: &account,
							Arn:     &arn,
							UserId:  &userID,
						}, nil
					},
				},
			},
			&STSIdentity{
				Account: "342563637373",
				ARN:     "ashgajsdhgajsdg",
				UserID:  "asjkdhkasdhaksd",
			},
			false,
		},
		{
			"Invaild/Error",
			args{
				stsInstance: &stsmock.STSAPIMock{
					GetCallerIdentityFunc: func(in1 *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
						return nil, xerrors.New("blah")
					},
				},
			},
			nil,
			true,
		},
		{
			"Invaild/awserrError",
			args{
				stsInstance: &stsmock.STSAPIMock{
					GetCallerIdentityFunc: func(in1 *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
						return nil, awserr.New("askjdhaksjhd", "Blah", xerrors.New("blah"))
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSTSIdentity(tt.args.stsInstance)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSTSIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSTSIdentity() = %v, want %v", got, tt.want)
			}
		})
	}
}
