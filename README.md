# GitHub Actions CLI (`ghactl`)

![GitHub Release (latest SemVer)](https://img.shields.io/github/v/release/action-stars/ghactl?logo=github&label=Release&sort=semver)
![Validate](https://github.com/action-stars/ghactl/actions/workflows/validate.yaml/badge.svg?branch=main)

`ghactl` is a CLI for GitHub Actions workflows. It exposes the same capabilities as the JavaScript [GitHub Actions Toolkit](https://github.com/actions/toolkit) as shell commands, so you can use them directly in `run:` steps without writing a custom action.

## Why ghactl?

The official toolkit packages only work in JavaScript/TypeScript actions. [Workflow commands](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions) exist for shell steps but the coverage is limited and the syntax is fiddly. `ghactl` is a single binary that gives shell steps access to tool-cache management, downloads, archive extraction, and version resolution.

## Features

- Add entries to the GitHub Actions PATH for subsequent workflow steps
- Install and cache tools from GitHub Releases
- Find, list, and cache tools in the runner tool cache with semver version matching
- Download files from URLs with automatic retries
- Extract `.tar`, `.tar.gz`, and `.zip` archives
- Check versions against semver constraints

## Installation

Installation instructions are not yet available.

## Commands

`ghactl` is organised into top-level commands, each covering a distinct area of GitHub Actions functionality.

> [!NOTE]
> More commands will be added over time.

| Command | Description                  |
| ------- | ---------------------------- |
| `path`  | Manage GitHub Actions PATH entries. |
| `tool`  | Manage GitHub runner tools.  |

### Global Flags

| Flag        | Description             |
| ----------- | ----------------------- |
| `--verbose` | Enable verbose output.  |
| `--version` | Print the version.      |
| `--help`    | Show help.              |

---

## `tool`

Manage GitHub runner tools: download, extract, cache, and check versions.

| Subcommand       | Description                                          |
| ---------------- | ---------------------------------------------------- |
| `cache get`      | Get the tool cache directory path.                   |
| `cache find`     | Find one or more cached tool versions.               |
| `cache add dir`  | Add a directory to the tool cache.                   |
| `cache add file` | Add a file to the tool cache.                        |
| `download`       | Download a tool to a temporary directory.            |
| `extract tar`    | Extract a tar archive to a temporary directory.      |
| `extract tgz`    | Extract a tar.gz archive to a temporary directory.   |
| `extract zip`    | Extract a zip archive to a temporary directory.      |
| `install`        | Install a tool from a source.                        |
| `version check`  | Check if a version matches a constraint.             |

---

## `path`

Manage GitHub Actions PATH entries.

| Subcommand  | Description        |
| ----------- | ------------------ |
| `add`       | Add a path entry.  |

---

### `path add`

Add a path entry.

This writes to the GitHub Actions `GITHUB_PATH` file to be used in future steps.

| Flag     | Required | Description          |
| -------- | -------- | -------------------- |
| `--path` | Yes      | Path entry to add.   |

```sh
ghactl path add --path "$HOME/.local/bin"
```

GitHub Actions step example:

```yaml
- name: Add local bin directory
  run: ghactl path add --path "${HOME}/.local/bin"
```

---

### `tool cache get`

Get the tool cache directory path.

```sh
ghactl tool cache get
```

---

### `tool cache find`

Find a specific cached tool version path, or all matching cached versions.

| Flag        | Required | Default        | Description                          |
| ----------- | -------- | -------------- | ------------------------------------ |
| `--name`    | Yes      |                | Name of the tool.                    |
| `--arch`    | No       | Runtime GOARCH | Architecture of the tool.            |
| `--version` | No       | `*` (any)      | Version spec to match.               |
| `--all`     | No       | `false`        | Return all matching cached versions. |

```sh
ghactl tool cache find --name my-tool --version "^1.0.0"
ghactl tool cache find --name my-tool --version "1.2.3" --arch arm64
ghactl tool cache find --name my-tool --all
ghactl tool cache find --name my-tool --all --version "^1.0.0"
```

---

### `tool cache add dir`

Add a directory to the tool cache.

| Flag       | Required | Default        | Description                 |
| ---------- | -------- | -------------- | --------------------------- |
| `--source` | Yes      |                | Source directory path.      |
| `--name`   | Yes      |                | Name of the tool.           |
| `--version`| Yes      |                | Version of the tool.        |
| `--arch`   | No       | Runtime GOARCH | Architecture of the tool.   |

```sh
ghactl tool cache add dir --source /tmp/extracted --name my-tool --version 1.2.3
```

---

### `tool cache add file`

Add a single file to the tool cache.

| Flag            | Required | Default        | Description                       |
| --------------- | -------- | -------------- | --------------------------------- |
| `--source`      | Yes      |                | Source file path.                 |
| `--name`        | Yes      |                | Name of the tool.                 |
| `--version`     | Yes      |                | Version of the tool.              |
| `--arch`        | No       | Runtime GOARCH | Architecture of the tool.         |
| `--target-name` | No       | Tool name      | Name to rename the source file to. |

```sh
ghactl tool cache add file --source /tmp/my-binary --name my-tool --version 1.0.0
ghactl tool cache add file --source /tmp/binary --name my-tool --version 1.0.0 --target-name custom-name
```

---

### `tool download`

Download a tool from a URL to a temporary directory. Outputs the path to the downloaded file.

| Flag    | Required | Description                    |
| ------- | -------- | ------------------------------ |
| `--url` | Yes      | URL to download the tool from. |

```sh
ghactl tool download --url https://example.com/tool-v1.0.0-linux-amd64.tar.gz
```

---

### `tool extract tar`

Extract a tar archive to a temporary directory.

| Flag     | Required | Description            |
| -------- | -------- | ---------------------- |
| `--path` | Yes      | Path to the tar archive. |

```sh
ghactl tool extract tar --path /tmp/tool.tar
```

---

### `tool extract tgz`

Extract a tar.gz archive to a temporary directory.

| Flag     | Required | Description                 |
| -------- | -------- | --------------------------- |
| `--path` | Yes      | Path to the tar.gz archive. |

```sh
ghactl tool extract tgz --path /tmp/tool.tar.gz
```

---

### `tool extract zip`

Extract a zip archive to a temporary directory.

| Flag     | Required | Description            |
| -------- | -------- | ---------------------- |
| `--path` | Yes      | Path to the zip archive. |

```sh
ghactl tool extract zip --path /tmp/tool.zip
```

---

### `tool install`

Install a tool from GitHub Releases and cache it in the GitHub runner tool cache.

| Flag            | Required | Default           | Description |
| --------------- | -------- | ----------------- | ----------- |
| `--owner`       | Yes      |                   | GitHub repository owner. |
| `--repo`        | Yes      |                   | GitHub repository name. |
| `--version`     | No       | `latest`          | Version input (`latest` or exact version/tag). |
| `--token`       | No       | `GITHUB_TOKEN`    | GitHub token used for API access. |
| `--name`        | No       | Value of `--repo` | Tool cache name. |
| `--arch`        | No       | Runtime GOARCH    | Tool architecture. |
| `--os`          | No       | Runtime GOOS      | Tool operating system. |
| `--pre-release` | No       | `false`           | Include pre-releases when resolving `latest`. |
| `--add-to-path` | No       | `true`            | Add the tool directory to PATH. |

```sh
ghactl tool install --owner cli --repo cli --version 2.94.0
```

---

### `tool version check`

Check if a version satisfies a semver constraint. Outputs `true` or `false`.

| Flag             | Required | Description                     |
| ---------------- | -------- | ------------------------------- |
| `--version`      | Yes      | Version to check.               |
| `--version-spec` | Yes      | Version constraint to check against. |

```sh
ghactl tool version check --version 1.2.3 --version-spec "^1.0.0"
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## Security

See [SECURITY.md](SECURITY.md) for reporting vulnerabilities.

## License

[MIT](LICENSE)
