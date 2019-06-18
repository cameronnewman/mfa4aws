package awssts

import (
	"os/user"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/spf13/afero"

	"mfa4aws/internal/pkg/awssts/mock/iammock"
	"mfa4aws/internal/pkg/awssts/mock/stsmock"

	"golang.org/x/xerrors"
)

func init() {
	appFs = afero.NewMemMapFs()

	//Known Path
	afero.WriteFile(appFs, "/knowntestfile.txt", []byte(`test`), 0644)
	afero.WriteFile(appFs, "/emptyknowntestfile.txt", []byte(nil), 0000)

	//$HOME/.aws/credentials
	user, _ := user.Current()
	path := filepath.Join(user.HomeDir, ".aws")
	credentialsFile := filepath.Join(path, "credentials")

	appFs.MkdirAll(path, 0755)
	afero.WriteFile(appFs, credentialsFile, []byte(`
	[default]
	aws_access_key_id = blahblah
	aws_secret_access_key = blahblah/blahblah`), 0644)
}

func TestGenerateSTSCredentials(t *testing.T) {
	type args struct {
		profile   string
		tokenCode string
	}
	tests := []struct {
		name    string
		args    args
		want    *AWSCredentials
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSTSCredentials(tt.args.profile, tt.args.tokenCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSTSCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSTSCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateProfile(t *testing.T) {
	type args struct {
		file    []byte
		profile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Valid/NoProfileNameDefined",
			args{
				file: []byte(`
				[default]
				aws_access_key_id = blahblah
				aws_secret_access_key = blahblah/blahblah`),
				profile: "",
			},
			false,
		},
		{
			"Valid/NonDefaultProfileNameDefined",
			args{
				file: []byte(`
				[candycrush]
				aws_access_key_id = blahblah
				aws_secret_access_key = blahblah/blahblah`),
				profile: "candycrush",
			},
			false,
		},
		{
			"Invalid/InvalidCredentialsFile",
			args{
				file: []byte(`
				-[default]
				_aws_access_key_id = blahblah
				a-ws_secret_access_key = blahblah/blahblah`),
				profile: "",
			},
			true,
		},
		{
			"Invalid/InvalidProfile",
			args{
				file:    []byte(""),
				profile: "blah",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateProfile(tt.args.file, tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("validateProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_openFile(t *testing.T) {

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"Invalid/NonExistentFile",
			args{
				path: "/some/unknown/path",
			},
			nil,
			true,
		},
		{
			"Valid/FileExists",
			args{
				path: "/knowntestfile.txt",
			},
			[]byte(`test`),
			false,
		},
		{
			"Valid/EmptyFileExists",
			args{
				path: "/emptyknowntestfile.txt",
			},
			[]byte(""),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("openFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("openFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Invalid/EmptyToken",
			args{
				token: "",
			},
			true,
		},
		{
			"Invalid/ShortToken",
			args{
				token: "2321",
			},
			true,
		},
		{
			"Valid/Token",
			args{
				token: "1234532",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateToken(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("validateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getIAMUserMFADevice(t *testing.T) {
	type args struct {
		iamInstance iamiface.IAMAPI
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Vaild/DeviceFound",
			args{
				iamInstance: &iammock.IAMAPIMock{
					ListMFADevicesFunc: func(in1 *iam.ListMFADevicesInput) (*iam.ListMFADevicesOutput, error) {
						sn := "shsjdyshe"

						output := &iam.ListMFADevicesOutput{
							MFADevices: []*iam.MFADevice{&iam.MFADevice{
								SerialNumber: &sn,
							}},
						}
						return output, nil
					},
				},
			},
			"shsjdyshe",
			false,
		},
		{
			"Invaild/awserrError",
			args{
				iamInstance: &iammock.IAMAPIMock{
					ListMFADevicesFunc: func(in1 *iam.ListMFADevicesInput) (*iam.ListMFADevicesOutput, error) {
						return nil, awserr.New("5000", "blah", xerrors.New("blah"))
					},
				},
			},
			"",
			true,
		},
		{
			"Invaild/Error",
			args{
				iamInstance: &iammock.IAMAPIMock{
					ListMFADevicesFunc: func(in1 *iam.ListMFADevicesInput) (*iam.ListMFADevicesOutput, error) {
						return nil, xerrors.New("blah")
					},
				},
			},
			"",
			true,
		},
		{
			"Invaild/NoDevices",
			args{
				iamInstance: &iammock.IAMAPIMock{
					ListMFADevicesFunc: func(in1 *iam.ListMFADevicesInput) (*iam.ListMFADevicesOutput, error) {

						output := &iam.ListMFADevicesOutput{
							MFADevices: nil,
						}
						return output, nil
					},
				},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIAMUserMFADevice(tt.args.iamInstance)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIAMUserMFADevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getIAMUserMFADevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
