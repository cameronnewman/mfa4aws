package cmd

import "testing"

func TestRootExecute(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Valid/ContainsVersionString",
			args{
				version: "1.0.1",
			},
		},
		{
			"Valid/ContainsNoVersionString",
			args{
				version: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute(tt.args.version)
		})
	}
}
