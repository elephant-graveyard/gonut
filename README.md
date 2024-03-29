# gonut [![License](https://img.shields.io/github/license/homeport/gonut.svg)](https://github.com/homeport/gonut/blob/main/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/homeport/gonut)](https://goreportcard.com/report/github.com/homeport/gonut) [![Build and Tests](https://github.com/homeport/gonut/workflows/Build%20and%20Tests/badge.svg)](https://github.com/homeport/gonut/actions?query=workflow%3A%22Build+and+Tests%22) [![Go Reference](https://pkg.go.dev/badge/github.com/homeport/gonut.svg)](https://pkg.go.dev/github.com/homeport/gonut) [![Release](https://img.shields.io/github/release/homeport/gonut.svg)](https://github.com/homeport/gonut/releases/latest) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/homeport/gonut)

![gonut](.docs/logo.png?raw=true "Gonut logo - molten chocolate covering a donut")

## Introducing gonut

`gonut` is a portable tool to help you verify whether you can push a sample app to a Cloud Foundry. It will push an app to Cloud Foundry and delete it afterwards. The apps are embedded into the `gonut` binary, so you just have to install `gonut` and you are set.

It is written in Go and uses [`pina-golada`](https://github.com/homeport/pina-golada) to include arbitrary sample app data in the final binary.

![gonut example](assets/images/gonut-example.gif?raw=true "Example of gonut pushing a sample app to Pivotal Cloud Foundry")

## How do I get started

Installation options are either using Homebrew or a convenience download script.

- On macOS systems, a Homebrew tap is available to install `gonut`:

  ```sh
  brew install homeport/tap/gonut
  ```

- Use a convenience script to download the latest release to install it in a suitable location on your local machine:

  ```sh
  curl -fsL http://ibm.biz/Bd2t2v | bash
  ```

## Contributing

We are happy to have other people contributing to the project. If you decide to do that, here's how to:

- get Go (`gonut` requires Go version 1.11 or greater)
- fork the project
- create a new branch
- make your changes
- open a PR.

Git commit messages should be meaningful and follow the rules nicely written down by [Chris Beams](https://chris.beams.io/posts/git-commit/):
> The seven rules of a great Git commit message
>
> 1. Separate subject from body with a blank line
> 1. Limit the subject line to 50 characters
> 1. Capitalize the subject line
> 1. Do not end the subject line with a period
> 1. Use the imperative mood in the subject line
> 1. Wrap the body at 72 characters
> 1. Use the body to explain what and why vs. how

### Running test cases and binaries generation

There are multiple make targets, but running `all` does everything you want in one call.

```sh
make all
```

### Test it with Linux on your macOS system

Best way is to use Docker to spin up a container:

```sh
docker run \
  --interactive \
  --tty \
  --rm \
  --volume $GOPATH/src/github.com/homeport/gonut:/go/src/github.com/homeport/gonut \
  --workdir /go/src/github.com/homeport/gonut \
  golang:1.19 /bin/bash
```

### Git pre-commit hooks

Add a pre-commit hook using this command in the repository directory:

```sh
cat <<EOS | cat > .git/hooks/pre-commit && chmod a+rx .git/hooks/pre-commit
#!/usr/bin/env bash

set -euo pipefail
make test

EOS
```
