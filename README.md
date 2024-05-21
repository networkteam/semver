# semver

[![GoDoc](https://godoc.org/github.com/networkteam/semver?status.svg)](https://godoc.org/github.com/networkteam/semver)
[![Build Status](https://github.com/networkteam/semver/workflows/Go/badge.svg)](https://github.com/networkteam/semver/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/networkteam/semver)](https://goreportcard.com/report/github.com/networkteam/semver)
[![codecov](https://codecov.io/gh/networkteam/semver/branch/main/graph/badge.svg?token=S8X8TMLQ9O)](https://codecov.io/gh/networkteam/semver)

A semver package for Go implementing the [SemVer 2.0 spec](https://semver.org/).

## Why?

We wanted an implementation that follows the spec as closely as possible.

## Install

```bash
go get github.com/networkteam/semver
```

## Example

```go
package main

import (
	"fmt"

	"github.com/networkteam/semver"
)

func main() {
	v, err := semver.ParseVersion("1.0.0-alpha.1+001")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Parsed Version", v.Major, v.Minor, v.Patch, v.PreRelease, v.Build)
		fmt.Println("String Representation:", v.String())
	}
}
```

## License

[MIT](./LICENSE)