# Multi-Git

A little Go command-line program that manages a list of git repos and runs git commands on all repos.

## Setting up Go and Installing the multi-git Package

### Prerequisites
- Go version 1.19 or higher
- Git command line tool

### Installation

- Install Go by following the instructions on the official [Go website](https://golang.org/doc/install).
- Open a command prompt and run `go version` to verify that Go is properly installed.
- Run the command `go install -v github.com/kphalgun/multi-git/cmd/multi-git@latest` to install the package.

### Command-line Arguments
It accepts two command-line arguments:

- command : the git command (wrap in double quotes for multi-arguments commands)
- ignore-errors: keeps going through the list of repos even the git command failed on some of them

### Environment variables
The list of repos is controlled by two the environment variables:

- MG_ROOT : the path to a root directory that contains all the repos
- MG_REPOS : the names of all managed repos under MG_ROOT

## Usage
- Open a command prompt.
- Run the command multi-git to use the package.

```
For example:

multi-git status

If you want to specify multiple flags to git surround them with quotes:

multi-git 'status --short'

It also requires the following environment variables defined:
MG_ROOT: root directory of target git repositories
MG_REPOS: list of repository names to operate on

Usage:
  $ multi-git [flags]

Flags:
      --ignore-errors   will continue executing the command for all repos if ignore-errors is true
                        otherwise it will stop execution when an error occurs
```

## Reference
1. https://github.com/the-gigi/multi-git
