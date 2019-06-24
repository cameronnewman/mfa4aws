package aws

import (
	"reflect"
	"testing"
)

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

func Test_createSessions(t *testing.T) {
	type args struct {
		path    string
		profile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Vaild/NoVariables",
			args{
				path:    "",
				profile: "",
			},
			false,
		},
		{
			"Invaild/InvaildPath",
			args{
				path:    "/shhss/ssjjss",
				profile: "",
			},
			true,
		},
		{
			"Invaild/InvaildPathInvaildProfile",
			args{
				path:    "/shhss/ssjjss",
				profile: "daskdjhaskdjhsd",
			},
			true,
		},
		{
			"Invaild/InvaildProfile",
			args{
				path:    "",
				profile: "akdjghakjsdhaksjdh",
			},
			true,
		},
		{
			"Vaild/VaildProfile",
			args{
				path:    "",
				profile: "default",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := createSession(tt.args.path, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("createSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
