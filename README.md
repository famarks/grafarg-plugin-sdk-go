# Grafarg Plugin SDK for Go

This SDK enables building [Grafarg](https://github.com/famarks/grafarg) backend plugins using Go.

[![License](https://img.shields.io/github/license/grafarg/grafarg-plugin-sdk-go)](LICENSE)
[![Go.dev](https://pkg.go.dev/badge/github.com/famarks/grafarg-plugin-sdk-go)](https://pkg.go.dev/github.com/famarks/grafarg-plugin-sdk-go?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/famarks/grafarg-plugin-sdk-go)](https://goreportcard.com/report/github.com/famarks/grafarg-plugin-sdk-go)
[![Circle CI](https://img.shields.io/circleci/build/gh/grafarg/grafarg-plugin-sdk-go/master)](https://circleci.com/gh/grafarg/grafarg-plugin-sdk-go?branch=master)

## Current state

This SDK is still in development. The protocol between the Grafarg server and the plugin SDK is considered stable, but we might introduce breaking changes in the SDK. This means that plugins using the older SDK should work with Grafarg, but might lose out on new features and capabilities that we introduce in the SDK.

## Contributing

If you're interested in contributing to this project:

- Start by reading the [Contributing guide](/CONTRIBUTING.md).
- Learn how to set up your local environment, in our [Developer guide](/contribute/developer-guide.md).

## License

[Apache 2.0 License](https://github.com/famarks/grafarg-plugin-sdk-go/blob/master/LICENSE)
