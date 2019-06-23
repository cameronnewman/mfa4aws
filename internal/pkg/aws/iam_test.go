package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"mfa4aws/internal/pkg/aws/mock/iammock"

	"golang.org/x/xerrors"
)

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
