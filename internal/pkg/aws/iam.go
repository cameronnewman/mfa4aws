package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

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
