/*
Simple CLI tool which enables you to login and retrieve [AWS](https://aws.amazon.com/) temporary credentials using for IAM users.
*/
package main

// Usage:
//   mfa4aws [flags]
//   mfa4aws [command]
//
// Available Commands:
//   help        Help about any command
//   shell      Generates AWS STS access keys for use on the shell by wrapping the result in eval
//
// Flags:
//   -h, --help             help for mfa4aws
//   -p, --profile string   AWS Profile name in $HOME/.aws/credentials (default "default")
//   -t, --token string     Current MFA value to use for STS generation
//
// Use "mfa4aws [command] --help" for more information about a command.
//
