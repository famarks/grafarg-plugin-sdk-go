# Developer guide

This guide helps you get started developing Grafarg Plugin SDK for Go.

## Tooling

Make sure you have the following tools installed before setting up your developer environment:

- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/) (see [go.mod](../go.mod#L3) for minimum required version)
- [Mage](https://magefile.org/)

## Building

We use [Mage](https://magefile.org/) as our primary tool for development related tasks like building and testing etc. It should be run from the root of this repository.

List available Mage targets that are available:

```bash
mage -l
```

You can use the `build` target to verify all code compiles. It  doesn't output any binary though.

```bash
mage -v build
```

The `-v` flag can be used to show verbose output when running Mage targets.

### Testing

```bash
mage test
```

### Linting

```bash
mage lint
```

### Generate Go code for Protobuf definitions

A prerequisite is to have [protoc](http://google.github.io/proto-lens/installing-protoc.html) installed and available in your path.

Next, you need to have [protoc-gen-go](https://github.com/golang/protobuf/tree/v1.3.4#installation) installed and available in your path. It's very important that you match the version specified of `github.com/golang/protobuf` in [go.mod](go.mod) file, as of this writing it's v1.3.4.

```
mage protobuf
```

### Changing `generic_*.go` files in the `data` package

Currently [genny](https://github.com/cheekybits/genny) is used for generating some go code. If you make changes to generic template files then `genny` needs to be installed, and then `mage dataGenerate`. Changed generated files should be committed with the change in the template files.

### Dependency management

We use Go modules for managing Go dependencies. After you've updated/modified modules dependencies, please run `go mod tidy` to cleanup dependencies.
