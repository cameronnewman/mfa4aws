package shell

import (
	"bytes"
	"testing"
)

func TestPrintVars(t *testing.T) {
	type args struct {
		vars []string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			"Invalid/EmptyCreds",
			args{
				vars: []string{},
			},
			"",
		},
		{
			"Valid/SimpleString",
			args{
				vars: []string{"Hi", "John"},
			},
			"Hi\nJohn\n",
		},
		{
			"Valid/ComplexString",
			args{
				vars: []string{"Hi", "世界"},
			},
			"Hi\n世界\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintVars(out, tt.args.vars)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("PrintVars() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
