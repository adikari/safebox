# 📦  SafeBox

SafeBox is a command line tool for managing secrets for your application. Currently it supports AWS Parameter Store and AWS Secrets Manager.

## Installation

SafeBox is available for many Linux distros and Windows.

```bash
# Via brew (OSX)
$ brew install adikari/taps/safebox

# Via curl
$ curl -sSL https://raw.githubusercontent.com/monebag/safebox/main/scripts/install.sh | sh

# Via npm
$ npm install @adikari/safebox

# Via yarn
$ yarn add @adikari/safebox
```

To install it directly find the right version for your machine in [releases](https://github.com/monebag/safebox/releases) page. Download and un-archive the files. Copy the `safebox` binary to the PATH or use it directly.

## Usage

1. Create a configuration file called `safebox.yml`.

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/monebag/safebox/main/schema.json
service: my-service
provider: ssm

config:
  defaults:
    DB_NAME: "database name updated"
    API_ENDPOINT: "http://some-endpoint-{{ .stage }}.com"

  prod:
    DB_NAME: "production db name"

  shared:
    SHARED_VARIABLE: "some shared config"

secret:
  defaults:
    API_KEY: "key of the api endpoint"
    DB_SECRET: "database secret"

  shared:
    SHARED_KEY: "shared key"
```

2. Use `safebox` CLI tool to deploy your configuration.

```bash
$ safebox deploy --stage <stage> --config path/to/safebox.yml --prompt="missing"
```

You can then run list command to view the pushed configurations.

The variables under
1. `defaults` is deployed with path prefix of `/<stage>/<service>`
1. `shared` is deployed with path prefix of `/<stage>/shared/`

### CLI Reference

Following are all options available in `safebox` CLI.

```bash
A Fast and Flexible secret manager built with love by adikari in Go.

Usage:
  safebox [flags]
  safebox [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  deploy      Deploys all configurations specified in config file
  export      Exports all configuration to a file
  help        Help about any command
  import      Imports all configuration from a file
  list        Lists all the configs available

Flags:
  -c, --config string   path to safebox configuration file (default "safebox.yml")
  -h, --help            help for safebox
  -s, --stage string    stage to deploy to (default "dev")
  -v, --version         version for safebox

Use "safebox [command] --help" for more information about a command.
```

### Using in scripts

```bash
#!/bin/bash

set -euo pipefail

echo "📦  deploying configs to ssm"
safebox deploy --stage $STAGE # ensures all configs are deployed. throws error if ay configs are missings

configs=$(safebox export --stage $STAGE)
CONFIG1=$(echo "$configs" | jq -r ".CONFIG1")
CONFIG2=$(echo "$configs" | jq -r '.CONFIG2')

echo $CONFIG1
echo $CONFIG2
```

### Generating dotenv files

This is quite handy when your build process or application requires configuration in a dotenv file. The command reads all your configs defined in `safebox.yml` and outputs the dotenv file.

```bash
safebox export --stage <stage> --format="dotenv" --output-file=".env"
```

### Replacing existing configuration

To replace the configuration simply update the value in the `safebox.yml` file and redeploy.
To replace the existing secrets run the following command

```bash
safebox deploy --stage <stage> --prompt="all"
```

This will display a prompt with the secret and its existing values. You can press enter to retain the old value for secrets that you don't want to update.
For the secret that you want to replace, remove the old value from the prompt then provide the new value.

### Deploy new configuration

To deploy the new configuration, simply add the new key value in `safebox.yml`
To deploy new secret value, run the following command

```bash
safebox deploy --stage <stage> --prompt="missing"
```

The missing flag will only prompt you for the new secrets.

### Configuration File Reference

Following is the configuration file will all possible options:

```yaml
service: my-service
provider: secrets-manager                     # ssm OR secrets-manager
prefix: "/custom/prefix/{{.stage}}/"          # Optional. Defaults to /<stage>/<service>/. Prefix all parameters. Does not apply for shared

stacks:                                       # Outputs from cloudformation stacks that needs to be interpolated.
  - some-cloudformation-stack

config:
  defaults:                                   # Default parameters. Can be overwritten in different environments.
    DB_NAME: my-database
    DB_HOST: 3200
    KEY_VALUE_SECRET: '{"hello": "world"}'    # JSON body can be passed when provider is secrets-manager. This will create key value secret
  production:                                 # If keys are deployed to production stage, its value will be overwritten by following
    DB_NAME: my-production-database
  shared:                                     # shared configuartions deployed under /<stage>/shared/ path
    DB_TABLE: "table-{{.stage}}"

secret:
  defaults:
    DB_PASSWORD: "secret database password"   # Value in quote is deployed as description of the ssm parameter.
```

**Variables available for interpolation**
- stage    - Stage used for deployment
- service  - Name of service as configured in the config file
- account  - AWS Account number
- region   - AWS Region

If using `stacks` then the outputs of that Cloudformation stack is also available for interpolation.

### Release

1. Update version number [npm/package.json](https://github.com/monebag/safebox/blob/main/npm/package.json).
2. Merge the changes to main branch.
2. Create a git tag that matches the same version number as npm package version.
3. Push the tag to github. Tag must follow semversion and prefixed with `v`. Eg. `v.1.2.3`.
4. Pushing the tag triggers github workflow that will automatically release new version.


### License

Feel free to use the code, it's released using the MIT license.
