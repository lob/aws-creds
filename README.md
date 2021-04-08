# aws-creds

A CLI tool to authenticate with Okta as the IdP to fetch AWS credentials.

## Installation

> Note: Building `aws-creds` requires `go >= 1.13`.

To install `aws-creds`:

1. Clone this repository
2. Run `make build`
3. Copy `./bin/aws-creds` to `/usr/local/bin`, or your preferred executables folder.

> Note: You can set the `BIN_DIR` and `VERSION` environment variables to override the build destination and version number.
>
> For example, building aws-creds into my `~/home/bin` directory:
> ```
> [aws-creds]$ VERSION=vX.Y.Z BIN_DIR=$HOME/bin make build
> [aws-creds]$ ~/bin/aws-creds -v
> aws-creds vX.Y.Z
> ```

## Usage

Once you have installed `aws-creds`, create a file `$HOME/.aws-creds/config` with content something like this:

```json
{
  "username": "<YOUR_OKTA_USER_NAME>", // Your Okta username
  "okta_host": "https://<YOUR_COMPANY>.okta.com", // Your organization's Okta URL
  "okta_app_path": "/home/amazon_aws/XXXXXXXXXXXXXXXXXXXX/###", // Okta app path for AWS
  "preferred_factor_type": "token:software:totp",
  "profiles": [
    {
      "name": "<PROFILE_NAME>", // One of your AWS profiles
      "role_arn": "arn:aws:iam:<ACCOUNT_ID>:role/<ROLE_NAME>" // The AWS Role to assume
    },
    // Other profiles
  ]
}
```

> Note: Remove all `//` comments in the final `.json` file.

Then in your `$HOME/.aws/config` make sure you have something like this:

```
[profile <PROFILE_NAME>]
region = <PREFERRED_AWS_REGION>
```

Then run

```
aws -p PROFILE_NAME
```

This may prompt you for your Okta password, Preferred multi-factor auth method, and then your one-time password.

`aws-creds` will populate your `$HOME/.aws/credentials` file with credentials from Okta.

Then you can run aws commands like:

```
aws sts get-caller-identity --profile PROFILE_NAME
```

Eventually these credentials will expire and you will need to re-run `aws -p PROFILE_NAME`.

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