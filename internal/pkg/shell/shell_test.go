package shell

import (
	"mfa4aws/internal/pkg/awssts"
	"reflect"
	"testing"
)

func TestGenerateSTSKeysForExport(t *testing.T) {
	type args struct {
		profile   string
		tokenCode string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenerateSTSKeysForExport(tt.args.profile, tt.args.tokenCode)
		})
	}
}

func Test_credentialsToEnvExport(t *testing.T) {
	type args struct {
		creds *awssts.AWSCredentials
	}
	tests := []struct {
		name        string
		args        args
		wantEnvVars []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEnvVars := credentialsToEnvExport(tt.args.creds); !reflect.DeepEqual(gotEnvVars, tt.wantEnvVars) {
				t.Errorf("credentialsToEnvExport() = %v, want %v", gotEnvVars, tt.wantEnvVars)
			}
		})
	}
}

func Test_printVars(t *testing.T) {
	type args struct {
		vars []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printVars(tt.args.vars)
		})
	}
}
