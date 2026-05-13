# dextools

`dextools` is a small Go module designed to keep useful tools together in one package. The root `dextools` package is the main entry point, and individual tools are added there as the module grows.

## Install

```bash
go get github.com/barluscuda/dextools
```

## Philosophy

Use `dextools` as the package you import first. Internally, the module contains focused subpackages, but the goal is to expose those tools through the root package so users get one clean entry point.

## Included Tools

This module currently includes:

- `dextools` as the unified entry package
- `envtools` for environment variable loading helpers
- `minecraft` for Minecraft developer tools

More tools can be added over time without changing the main idea: `dextools` is the umbrella package.

## Current Root API

The root package currently exposes these tools today:

- `Env()` for environment variable helpers
- `MinecraftBE()` for the current Bedrock helper

## Using Subpackages

The root package is the preferred way to use `dextools`. Subpackages are still available directly, but they support the broader goal of keeping tools collected under one package.

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

Example:

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

For package-level behavior and examples, see [envtools/DOCS.md](./envtools/DOCS.md).

## minecraft

The `minecraft` package is for Minecraft developer tools.

Right now, `MinecraftBE()` returns the current Bedrock helper from that package. This is only one part of the broader `minecraft` package.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/barluscuda/dextools"
)

func main() {
	client := dextools.MinecraftBE()
	client.Timeout = 5 * time.Second

	status, err := client.ServerStatus("play.nethergames.org", 19132)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MOTD:", status.MOTD)
	fmt.Println("Version:", status.Version)
	fmt.Printf("Players: %d/%d\n", status.PlayersOnline, status.PlayersMax)
	fmt.Println("Ping:", status.Ping)
}
```

Current response fields from the Bedrock status helper include:

- `Ping`
- `MOTD`
- `Protocol`
- `Version`
- `PlayersOnline`
- `PlayersMax`
- `ServerID`
- `LevelName`
- `GameMode`
- `GameModeNumeric`
- `IPv4Port`
- `IPv6Port`

For package-level behavior and examples, see [minecraft/DOCS.md](./minecraft/DOCS.md).

## Testing

Run all tests with:

```bash
go test ./...
```

GitHub Actions is configured to run the test suite on `push` and `pull_request`.
