# taggit

A CLI tool to publish semver releases from a git repository.

## Overview

`taggit` is a command line tool written in Go that helps manage semantic versioning
tags in a git repository. It provides commands for listing existing tags and creating
new version tags at patch, minor, and major levels.

## Getting Started

`taggit` can be installed via `go install` or downloaded from [Releases](https://github.com/shoenig/taggit/releases).

#### Install

```shell
go install github.com/shoenig/taggit@latest
```

## Commands

The `taggit` command line tool provides commands for managing version tags.

```
NAME:
  taggit

USAGE:
  taggit  [global options] [command [command options]] [arguments...]

VERSION:
  development

DESCRIPTION:
  Publish new versions of Go modules.

COMMANDS:
  list  - List tagged versions.
  zero  - Create initial v0.0.0 tag
  patch - Create an incremented patch version
  minor - Create an incremented minor version
  major - Create an incremented major version

GLOBALS:
--version/-V   boolean - print version information
--help/-h      boolean - print help message
```

#### list

The `list` command displays all semver-compatible tags in the repository, organized by base version.

```shell
$ taggit list
v0.1.0 |= v0.1.0 v0.1.0-alpha1
v0.2.0 |= v0.2.0-rc1 v0.2.0-r1+linux v0.2.0-r1+darwin
```

#### zero

The `zero` command creates the initial `v0.0.0` tag. This must be done before any other
version commands can be used.

```shell
$ taggit zero
taggit: created tag v0.0.0
```

#### patch

The `patch` command increments the patch level of the current latest version.

```shell
$ taggit patch
taggit: created tag v0.0.1
```

To create a pre-release version, provide the pre-release identifier as an argument:

```shell
$ taggit patch beta1
taggit: created tag v0.0.2-beta1
```

To add build metadata, use the `-m` or `--meta` flag:

```shell
$ taggit patch -m build123
taggit: created tag v0.0.2+build123
```

#### minor

The `minor` command increments the minor level of the current latest version.

```shell
$ taggit minor
taggit: created tag v0.1.0
```

Pre-release and metadata options work the same as with `patch`:

```shell
$ taggit minor alpha1
taggit: created tag v0.2.0-alpha1

$ taggit minor --meta=linux
taggit: created tag v0.2.0+linux
```

#### major

The `major` command increments the major level of the current latest version.

```shell
$ taggit major
taggit: created tag v1.0.0
```

## Global Options

- `-V, --version` - Print version information
- `-h, --help` - Print help message

## Contributing

The `github.com/shoenig/taggit` module is always improving. For bug fixes and new
features please file an issue.

## License

The `github.com/shoenig/taggit` module is open source under the [BSD-3-Clause](LICENSE) license.