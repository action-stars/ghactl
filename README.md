# GitHub Actions CLI (`ghactl`)

`ghactl` is a CLI for interacting with GitHub Actions from within a GitHub Actions workflow. `ghactl` is intended to be a friendlier way to make use of the [workflow commands for GitHub Actions](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions) as well as adding support for some capabilities currently only available in the [GitHub Actions Toolkit](https://github.com/actions/toolkit).

## Usage

```shell
> ghactl --help

NAME:
   ghactl - CLI to interact with GitHub Actions

USAGE:
   ghactl [global options] command [command options]

VERSION:
   (devel)

COMMANDS:
   env      Manage workflow environment variables
   matcher  Manage problem matchers
   message  Manage workflow log messages
   output   Manage workflow step outputs
   path     Manage workflow PATH configuration
   secret   Manage workflow secrets
   summary  Manage workflow summary messages
   temp     Get the runner temporary directory
   tool     Manage GitHub runner tools
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
