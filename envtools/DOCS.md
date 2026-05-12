# envtools

`envtools` provides small helpers for reading environment variables and converting them into common Go types.

## Quick Start

Use the package directly:

```go
package main

import (
	"fmt"
	"time"

	"github.com/barluscuda/dextools/envtools"
)

func main() {
	env := envtools.EnvTools{}

	host := env.LoadString("APP_HOST", "127.0.0.1")
	port := env.LoadPort("APP_PORT", 8080)
	debug := env.LoadBool("APP_DEBUG", false)
	timeout := env.LoadDuration("APP_TIMEOUT", 5*time.Second)

	fmt.Println(host, port, debug, timeout)
}
```

## Behavior

Every type has two forms:

- `Load*` returns the parsed value when possible.
- `Load*` returns the first fallback argument when the variable is missing or invalid.
- `Load*` returns the zero value when the variable is missing or invalid and no fallback is provided.
- `MustLoad*` panics when the variable is missing or invalid.

## Available Helpers

The package currently supports:

- `LoadString` / `MustLoadString`
- `LoadInt` / `MustLoadInt`
- `LoadInt64` / `MustLoadInt64`
- `LoadBool` / `MustLoadBool`
- `LoadFloat64` / `MustLoadFloat64`
- `LoadDuration` / `MustLoadDuration`
- `LoadTime` / `MustLoadTime`
- `LoadStringSlice` / `MustLoadStringSlice`
- `LoadBase64` / `MustLoadBase64`
- `LoadPort` / `MustLoadPort`

## Parsing Rules

### Strings

`LoadString` reads the value as-is.

```go
name := env.LoadString("APP_NAME", "demo")
```

### Integers

`LoadInt` uses `strconv.Atoi`.

```go
workers := env.LoadInt("APP_WORKERS", 4)
```

`LoadInt64` uses `strconv.ParseInt(..., 10, 64)`.

```go
limit := env.LoadInt64("APP_LIMIT", 1000)
```

### Booleans

`LoadBool` uses `strconv.ParseBool`, so values such as `true`, `false`, `1`, and `0` are accepted.

```go
debug := env.LoadBool("APP_DEBUG", false)
```

### Floats

`LoadFloat64` uses `strconv.ParseFloat(..., 64)`.

```go
ratio := env.LoadFloat64("APP_RATIO", 0.5)
```

### Durations

`LoadDuration` uses `time.ParseDuration`.

```go
timeout := env.LoadDuration("APP_TIMEOUT", 30*time.Second)
```

Example environment values:

- `500ms`
- `5s`
- `2m`
- `1h30m`

### Time

`LoadTime` expects RFC3339 format.

```go
startedAt := env.LoadTime("APP_STARTED_AT")
```

Example environment value:

```text
2026-05-12T10:30:00Z
```

### String Slices

`LoadStringSlice` splits on commas and trims spaces around each item.

```go
hosts := env.LoadStringSlice("APP_HOSTS", []string{"localhost"})
```

Example environment value:

```text
api-1.internal, api-2.internal, api-3.internal
```

Result:

```go
[]string{"api-1.internal", "api-2.internal", "api-3.internal"}
```

### Base64

`LoadBase64` expects standard base64 encoding and returns raw bytes.

```go
secret := env.LoadBase64("APP_SECRET")
```

### Port

`LoadPort` first parses the value as an integer, then enforces the valid TCP/UDP port range `1..65535`.

```go
port := env.LoadPort("APP_PORT", 8080)
```

If the parsed port is outside that range:

- `LoadPort` returns the fallback when provided.
- `LoadPort` returns `0` when no fallback is provided.
- `MustLoadPort` panics.

## Panic Behavior

`MustLoad*` helpers panic with messages like:

```text
envtools: required environment variable "APP_PORT" is missing
```

or:

```text
envtools: required environment variable "APP_PORT" is invalid (strconv.Atoi: parsing "abc": invalid syntax)
```

`MustLoadPort` can also panic when the number is parsed successfully but falls outside the valid range.

## Testing

Unit tests for the loader helpers live in [load_test.go](./load_test.go).
