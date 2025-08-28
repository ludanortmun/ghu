# GHU - Github Utilities

`ghu` is a collection of CLI tools for interacting with GitHub repositories.
Written in Go and leveraging the GitHub API and Google's [go-github library](https://github.com/google/go-github).


# Requirements

- Go 1.25 or higher

# Installation

To install the tool, run:

```bash
go install github.com/ludanortmun/ghu@latest
```

# Commands

## clonedir

Downloads a directory or file from a GitHub repository.

### Usage

To download a directory from a GitHub repository, use the following command:

```bash
ghu clonedir <github-url> [destination-path]
```

where:
- `<github-url>` is the URL of the GitHub repository and the path to the directory you want to download.
- `[destination-path]` is an optional argument specifying where to save the downloaded directory. If not provided, it defaults to the current directory.

## serve

Creates a local web server that serves files directly from the specified GitHub repository. Allows you to access and view static files hosted on GitHub through a web browser without needing to clone the repository.

### Usage
To serve static files from a GitHub repository, use the following command:

```bash
ghu serve <github-url> [flags]
```

Available flags:
- `--port` or `p`: Specify the port to serve on (default is 8080).

## auth

Manages GitHub authentication tokens for accessing private repositories. Other commands will use this token if available.
It uses the OS keyring to securely store and retrieve the token.

### Usage
To set or update your GitHub authentication token, use the following command:

```bash
ghu auth set-token
```

You will then be prompted to enter your GitHub token. Note that stdout will be disabled while entering the token for security reasons.

You can also clear the stored token with:

```bash
ghu auth clear-token
```