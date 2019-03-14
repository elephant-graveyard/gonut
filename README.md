# gonut

[![License](https://img.shields.io/github/license/homeport/gonut.svg)](https://github.com/homeport/gonut/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/homeport/gonut)](https://goreportcard.com/report/github.com/homeport/gonut)
[![Build Status](https://travis-ci.org/homeport/gonut.svg?branch=master)](https://travis-ci.org/homeport/gonut)
[![GoDoc](https://godoc.org/github.com/homeport/gonut?status.svg)](https://godoc.org/github.com/homeport/gonut)
[![Release](https://img.shields.io/github/release/homeport/gonut.svg)](https://github.com/homeport/gonut/releases/latest)

## Introducing gonut

`gonut` is a portable tool to help you verify whether you can push a sample app to a Cloud Foundry. It will push an app to Cloud Foundry and delete it afterwards. The apps are embedded into the `gonut` binary, so you just have to install `gonut` and you are set.

It is written in Golang and uses [`pina-golada`](https://github.com/homeport/pina-golada) to include arbitrary sample app data in the final binary.

![gonut example](assets/images/gonut-example.gif?raw=true "Example of gonut pushing a sample app to Pivotal Cloud Foundry")

_This project is work in progress._

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
  golang:1.11 /bin/bash
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
