# dextools

`dextools` is a small Go utility module that exposes focused helper packages through a simple root API.

## Install

```bash
go get github.com/barluscuda/dextools
```

## Current Root API

The root package currently exposes:

- `Env()` for environment variable helpers

## Quick Start

```go
package main

import (
	"fmt"
	"time"

	"github.com/barluscuda/dextools"
)

func main() {
	host := dextools.Env().LoadString("APP_HOST", "127.0.0.1")
	port := dextools.Env().LoadPort("APP_PORT", 8080)
	timeout := dextools.Env().LoadDuration("APP_TIMEOUT", 5*time.Second)

	fmt.Println(host, port, timeout)
}
```

## envtools

`Env()` returns an `envtools.EnvTools` value, which provides helpers for loading:

- strings
- integers
- int64 values
- booleans
- float64 values
- durations
- RFC3339 timestamps
- comma-separated string slices
- base64-encoded bytes
- validated ports

For package-level behavior and examples, see [envtools/DOCS.md](./envtools/DOCS.md).

## Testing

Run all tests with:

```bash
go test ./...
```

GitHub Actions is configured to run the test suite on `push` and `pull_request`.

Some live integration tests read repository secrets through environment variables:

- `MINECRAFT_MCSERVER_DEMO_IP`
- `MINECRAFT_MCSERVER_DEMO_PORT`
- `WENOWA_TOKEN`
- `WENOWA_DEMO_PHONE_NUMBER`
- `WENOWA_DEMO_HEADER`
- `WENOWA_DEMO_MESSAGE`
