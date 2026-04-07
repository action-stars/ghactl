# Contributing to ghactl

## Prerequisites

You need the following installed:

- [Go](https://go.dev/dl/) 1.26 or later
- [golangci-lint](https://golangci-lint.run/welcome/install/)
- [just](https://github.com/casey/just#installation) (command runner)
- [GoReleaser](https://goreleaser.com/install/) (optional, for release builds)

## Getting Started

1. Fork and clone the repository:

   ```sh
   git clone https://github.com/action-stars/ghactl.git
   cd ghactl
   ```

2. Run setup to tidy dependencies and create the output directory:

   ```sh
   just setup
   ```

## Development Commands

Common tasks are run through `just`:

| Command          | Description                              |
| ---------------- | ---------------------------------------- |
| `just setup`     | Tidy modules and create the dist folder  |
| `just tidy`      | Run `go mod tidy`                        |
| `just fmt`       | Format code with `golangci-lint fmt`     |
| `just lint`      | Lint and auto-fix with `golangci-lint`   |
| `just test`      | Run all tests with coverage              |
| `just build`     | Build the binary to `./dist`             |
| `just build-all` | Build all targets with GoReleaser        |

## Making Changes

1. Create a branch from `main`:

   ```sh
   git checkout -b my-feature
   ```

2. Make your changes. Keep commits focused.

3. Run linting and tests before pushing:

   ```sh
   just lint
   just test
   ```

4. Push your branch and open a pull request against `main`.

## Pull Request Guidelines

- One logical change per PR.
- Describe what the change does and why.
- All CI checks must pass.
- Add or update tests for new or changed behaviour.

## Reporting Issues

Found a bug or have a suggestion? [Open an issue](https://github.com/action-stars/ghactl/issues/new) with steps to reproduce and your environment details.
