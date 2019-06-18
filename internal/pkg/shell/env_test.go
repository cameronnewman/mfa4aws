package shell

import (
	"mfa4aws/internal/pkg/awssts"
	"reflect"
	"testing"
	"time"
)

func TestBuildEnvVars(t *testing.T) {
	type args struct {
		creds *awssts.AWSCredentials
	}
	tests := []struct {
		name        string
		args        args
		wantEnvVars []string
	}{
		{
			"Invalid/EmptyCreds",
			args{
				creds: &awssts.AWSCredentials{},
			},
			[]string{"export AWS_ACCESS_KEY_ID=", "export AWS_SECRET_ACCESS_KEY=", "export AWS_SESSION_TOKEN=", "export AWS_SECURITY_TOKEN=", "export X_PRINCIPAL_ARN=", "export EXPIRES=0s"},
		},
		{
			"Valid/Creds",
			args{
				creds: &awssts.AWSCredentials{
					AWSAccessKeyID:     "AHIAACNB4F5KCDQXSGYW4",
					AWSSecretAccessKey: "Xoy7ogSQXyTyZI3Oqv8JdAkk1PsbSYzt/vqQ1v+9",
					AWSSessionToken:    "FQoGZXIvYshgsSJHIOSLKj6nr0FOKIuOP68yKRKvPp3nj9MyaPcvN8PApmWd3yKuTJWf+u8hPmiDGIHAgDu5h+mVTdKL6B/gheTIjsqty9yn3it/2OoJNIhNfIPANfLwHnCSror+GDmS2Y/vZLjAThX0KKaM0/YcmUokHFMNrN+mAX8G21uAs0MUS4e5qzupfskjhskjhsk89797wZROPTk43ZharJLNf59hGVXnqHFwkxNatt/lKJH+pL0xScBr64qEb2ZaKOPonegF",
					AWSSecurityToken:   "FQoGZXIvYshgsSJHIOSLKj6nr0FOKIuOP68yKRKvPp3nj9MyaPcvN8PApmWd3yKuTJWf+u8hPmiDGIHAgDu5h+mVTdKL6B/gheTIjsqty9yn3it/2OoJNIhNfIPANfLwHnCSror+GDmS2Y/vZLjAThX0KKaM0/YcmUokHFMNrN+mAX8G21uAs0MUS4e5qzupfskjhskjhsk89797wZROPTk43ZharJLNf59hGVXnqHFwkxNatt/lKJH+pL0xScBr64qEb2ZaKOPonegF",
					PrincipalARN:       "162171167783:user/johnsmith",
					Expires:            time.Duration(1000000),
				},
			},
			[]string{"export AWS_ACCESS_KEY_ID=AHIAACNB4F5KCDQXSGYW4", "export AWS_SECRET_ACCESS_KEY=Xoy7ogSQXyTyZI3Oqv8JdAkk1PsbSYzt/vqQ1v+9", "export AWS_SESSION_TOKEN=FQoGZXIvYshgsSJHIOSLKj6nr0FOKIuOP68yKRKvPp3nj9MyaPcvN8PApmWd3yKuTJWf+u8hPmiDGIHAgDu5h+mVTdKL6B/gheTIjsqty9yn3it/2OoJNIhNfIPANfLwHnCSror+GDmS2Y/vZLjAThX0KKaM0/YcmUokHFMNrN+mAX8G21uAs0MUS4e5qzupfskjhskjhsk89797wZROPTk43ZharJLNf59hGVXnqHFwkxNatt/lKJH+pL0xScBr64qEb2ZaKOPonegF", "export AWS_SECURITY_TOKEN=FQoGZXIvYshgsSJHIOSLKj6nr0FOKIuOP68yKRKvPp3nj9MyaPcvN8PApmWd3yKuTJWf+u8hPmiDGIHAgDu5h+mVTdKL6B/gheTIjsqty9yn3it/2OoJNIhNfIPANfLwHnCSror+GDmS2Y/vZLjAThX0KKaM0/YcmUokHFMNrN+mAX8G21uAs0MUS4e5qzupfskjhskjhsk89797wZROPTk43ZharJLNf59hGVXnqHFwkxNatt/lKJH+pL0xScBr64qEb2ZaKOPonegF", "export X_PRINCIPAL_ARN=162171167783:user/johnsmith", "export EXPIRES=1ms"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEnvVars := BuildEnvVars(tt.args.creds); !reflect.DeepEqual(gotEnvVars, tt.wantEnvVars) {
				t.Errorf("BuildEnvVars() = %v, want %v", gotEnvVars, tt.wantEnvVars)
			}
		})
	}
}
