# aws-creds

A CLI tool to authenticate with Okta as the IdP to fetch AWS credentials.

## Installation

> Note: Building `aws-creds` requires `go >= 1.13`.

To install `aws-creds`:

1. Clone this repository.
2. Run `make build`.
3. Copy `./bin/aws-creds` to `/usr/local/bin`, or your preferred executables folder.

> Note: You can set the `BIN_DIR` and `VERSION` environment variables to override the build destination and version number.
>
> For example, building aws-creds directly into my `~/home/bin` directory:
> ```
> [aws-creds]$ VERSION=vX.Y.Z BIN_DIR=$HOME/bin make build
> [aws-creds]$ ~/bin/aws-creds -v
> aws-creds vX.Y.Z
> ```

## Usage

> Note: This assumes you have already installed the [AWS CLI](https://aws.amazon.com/cli/).

Once you have installed `aws-creds`, run `aws-creds configure`.
This will prompt you for your okta username, okta host + okta aws app path, and information about your aws profiles + role ARNs.

```
$ aws-creds configure
Configuring global settings...
Okta username: yourname
Okta AWS Embed Link (e.g. https://example.okta.com/home/amazon_aws/0oa54k1gk2ukOJ9nGDt7/252): https://example.okta.com/home/amazon_aws/0oa54k1gk2ukOJ9nGDt7/123

Configuring profile settings...
Profile name: test
Role ARN (e.g. arn:aws:iam::123456789001:role/EngineeringRole): arn:aws:iam::123456789001:role/SomeRole
Do you want to configure more profiles? [y/N]:
```

Once complete, `aws-creds` will create the file `$HOME/.aws-creds/config` with this information.

To fetch credentials, run `aws-creds -p`:

```
$ aws-creds -p $PROFILE_NAME
```

This may prompt you for your Okta password, preferred multi-factor auth method, and then your one-time password.

Please note that push-based authentication methods are not supported by this tool. You may need to add an additional verification method to your Okta account that uses an OTP. Google Authenticator and Okta Verify both work.

`aws-creds` will then populate your `$HOME/.aws/credentials` file with credentials from Okta.

Finally pass `--profile=$PROFILE_NAME`, or set `AWS_PROFILE=$PROFILE_NAME`, when running `aws` commands.
For example:

```
$ aws s3 ls --profile $PROFILE_NAME

$ AWS_PROFILE=sandbox aws s3 ls      # Equivalent

$ export AWS_PROFILE=sandbox         # Also
$ aws s3 ls                          # Equivalent
```

> Note: these credentials expire every hour, so you will need to re-run `aws-creds -p $PROFILE_NAME` periodically.

### Headless Linux / WSL Users

This utility uses a keyring library that is incompatible with certain linux configurations. To disable keyring caching, modify the `~/.aws-creds/config` file to include `"enable_keyring" : false,` after setting up the tool. 

## Building

`aws-creds` has a Makefile with helper commands:

make command | description
--- | ---
`make build` | Builds executable. Set `VERSION` and `BIN` environment variables to override defaults.
`make clean` | Clean `BIN` folder and go cache.
`make test` | Runs unit tests and generates coverage report.
`make install`| Installs go dependencies.
`make lint` | Lints the codebase.
`make release` | Creates and pushes a git tag.
`make setup` | Sets up linting and changelog tools.
`make html` | Generates test coverage report.
`make enforce` | Enforces test coverage.
