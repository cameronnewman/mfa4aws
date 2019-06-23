package aws

import (
	"os/user"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"mfa4aws/internal/pkg/aws/mock/iammock"

	"github.com/spf13/afero"

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
		want    *Credentials
		wantErr bool
	}{
		{
			"Invaild/NoVariables",
			args{
				profile:   "",
				tokenCode: "",
			},
			nil,
			true,
		},
		{
			"Invaild/InvaildProfile",
			args{
				profile:   "akdjghakjsdhaksjdh",
				tokenCode: "",
			},
			nil,
			true,
		},
		{
			"Invaild/InvaildTokenCode",
			args{
				profile:   "",
				tokenCode: "123",
			},
			nil,
			true,
		},
		{
			"Vaild/EmptyProfileAndVaildTokenCode",
			args{
				profile:   "",
				tokenCode: "1234564",
			},
			nil,
			true,
		},
		{
			"Vaild/VaildProfileAndVaildTokenCode",
			args{
				profile:   "default",
				tokenCode: "1234564",
			},
			nil,
			true,
		},
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
