package awssts

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
)

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

func Test_checkProfile(t *testing.T) {
	type args struct {
		profile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkProfile(tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("checkProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkToken(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("checkToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createNewSession(t *testing.T) {
	type args struct {
		profile string
	}
	tests := []struct {
		name    string
		args    args
		want    *session.Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createNewSession(tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("createNewSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createNewSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
