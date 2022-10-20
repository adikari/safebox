# 📦  SafeBox

SafeBox is a command line tool for managing secrets for your application. Currently it supports AWS Parameter Store.

## Installation

SafeBox is available for many Linux distros and Windows.

```bash
# Via brew (OSX)
$ brew install adikari/taps/safebox

# Via curl
$ curl -sSL https://raw.githubusercontent.com/adikari/safebox/main/scripts/install.sh | sh

# Via npm
$ npm install @adikari/safebox

# Via yarn
$ yarn add @adikari/safebox
```

To install it directly find the right version for your machine in [releases](https://github.com/adikari/safebox/releases) page. Download and un-archive the files. Copy the `safebox` binary to the PATH or use it directly.

## Usage

1. Create a configuration file called `safebox.yml`.

```yaml
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
$ safebox deploy --stage <stage> --config path/to/safebox.yml --prompt missing
```

You can then run list command to view the pushed configurations.

The variables under
1. `defaults` is deployed with path prefix of `/<stage>/<service>`
1. `shared` is deployed with path prefix of `/shared/`

### Config File

Following is the configuration file will all possible options:

```yaml
service: my-service
provider: ssm                                 # Only supports ssm for now.

stacks:                                       # Outputs from cloudformation stacks that needs to be interpolated.
  - some-cloudformation-stack

config:
  defaults:                                   # Default parameters. Can be overwritten in different environments.
    DB_NAME: my-database
    DB_HOST: 3200
  production:                                 # If keys are deployed to production stage, its value will be overwritten by following
    DB_NAME: my-production-database
  shared:                                     # shared configuartions deployed under /shared/ path
    DB_TABLE: "table-{{.stage}}"

secret:
  defaults:
    DB_PASSWORD: "secret database password"   # Value in quote is deployed as description of the ssm parameter.
```

### CLI

Following is all options available in `safebox` CLI.

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
yarn safebox deploy --stage $STAGE # ensures all configs are deployed. throws error if ay configs are missings

configs=$(yarn safebox export --stage $STAGE)
CONFIG1=$(echo "$configs" | jq -r ".CONFIG1")
CONFIG2=$(echo "$configs" | jq -r '.CONFIG2')

echo $CONFIG1
echo $CONFIG2

```

### Release

1. Update version number [npm/package.json](https://github.com/adikari/safebox/blob/main/npm/package.json).
2. Merge the changes to main branch.
2. Create a git tag that matches the same version number as npm package version.
3. Push the tag to github. Tag must follow semversion and prefixed with `v`. Eg. `v.1.2.3`.
4. Pushing the tag triggers github workflow that will automatically release new version.


### License

Feel free to use the code, it's released using the MIT license.
