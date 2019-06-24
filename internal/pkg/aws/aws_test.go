package aws

import (
	"os/user"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"
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
