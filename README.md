# mfa4aws  

[![Build][1]][2]
[![GoDoc][3]][4]
[![Go Report Card][5]][6]
[![FOSSA Status][9]][10]

[1]: https://github.com/cameronnewman/mfa4aws/workflows/pipeline/badge.svg
[2]: https://github.com/cameronnewman/mfa4aws/actions
[3]: https://godoc.org/github.com/cameronnewman/mfa4aws?status.svg
[4]: https://godoc.org/github.com/cameronnewman/mfa4aws
[5]: https://goreportcard.com/badge/github.com/cameronnewman/mfa4aws
[6]: https://goreportcard.com/report/github.com/cameronnewman/mfa4aws
[9]: https://app.fossa.io/api/projects/git%2Bgithub.com%2Fcameronnewman%2Fmfa4aws.svg?type=shield
[10]: https://app.fossa.io/projects/git%2Bgithub.com%2Fcameronnewman%2Fmfa4aws?ref=badge_shield

Simple CLI tool which enables you to login and retrieve [AWS](https://aws.amazon.com/) temporary credentials  for IAM users.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Requirements](#requirements)
- [Install](#install)
    - [OSX](#osx)
- [Usage](#usage)
    - [`mfa4aws shell`](#mfa4aws-shell)
- [Example](#example)
- [Building](#building)
- [Environment vars](#environment-vars)

## Requirements

* Access Key and Secret Key stores $HOME/.aws/credentials
* AWS IAM User account

## Install

### OSX

TBA

## Usage

```
Usage:
  shell [command]

Available Commands:
  help        Help about any command
  shell       Generates AWS STS access keys for use on the shell by wrapping the result in eval
  version     display release version

Flags:
  -h, --help             help for shell
  -p, --profile string   AWS Profile name in $HOME/.aws/credentials (default "default")
  -t, --token string     Current MFA value to use for STS generation

Use "shell [command] --help" for more information about a command.
```


### `mfa4aws shell`

If the `shell` sub-command is called, `mfa4aws` will output the following temporary security credentials:
```
export AWS_ACCESS_KEY_ID=DDFHAFG....UOCA
export AWS_SECRET_ACCESS_KEY="JSKA...HJ2F
export AWS_SESSION_TOKEN=ZQ...1VVQ==
export AWS_SECURITY_TOKEN=ZQ...1VVQ==
export X_PRINCIPAL_ARN=arn:aws:iam::3678236812376:user/johnsmith
```

If you use `eval $(mfa4aws shell)` frequently, you may want to create a alias for it:

zsh:
```
alias m4a="function(){eval $( $(command mfa4aws) shell --token=$@);}"
```

bash:
```
function m4a { eval $( $(which mfa4aws) shell --token=$@); }
```


## Building

```
make build
```

## Environment vars

The exec sub command will export the following environment variables.

* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY
* AWS_SESSION_TOKEN
* AWS_SECURITY_TOKEN
* X_PRINCIPAL_ARN

# License

This code is Copyright (c) 2018 [Cameron Newman](https://cameron.newman.io) and released under the MIT license. All rights not explicitly granted in the MIT license are reserved. See the included LICENSE.md file for more details.
