# mockingio

[![CI](https://github.com/mockingio/mockingio/actions/workflows/main.yml/badge.svg)](https://github.com/mockingio/mockingio/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/mockingio/mockingio/branch/main/graph/badge.svg?token=0AXGI7UR85)](https://codecov.io/gh/mockingio/mockingio)
[![Docker Repository](https://img.shields.io/docker/pulls/mockingio/mockingio)](https://hub.docker.com/r/mockingio/mockingio)
[![Github Release](https://img.shields.io/github/v/release/mockingio/mockingio)](https://github.com/mockingio/mockingio/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/mockingio/mockingio)](https://goreportcard.com/report/github.com/mockingio/mockingio)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


## Installation

### Go install

```go
go install github.com/mockingio/mockingio@latest
```

### Homebrew

```shell
brew tap mockingio/mockingio-tap
brew install mockingio/mockingio-tap/mockingio
```

### Docker hub

```shell
docker pull mockingio/mockingio

docker run -ti mockingio/mockingio --version
```

## Usage

### CLI

`mockingio start --filename example/mock.yml`
