# wenova

`wenova` provides Go helpers for calling Wenova microservice APIs from the `dextools` module.

## Quick Start

Use the package directly:

```go
package main

import (
	"context"
	"fmt"

	"github.com/barluscuda/dextools/wenova"
)

func main() {
	ctx := context.Background()
	client := wenova.Wenova{
		BaseUrl: "https://apimicroservices.wenova.fun",
	}

	smsResult, err := client.SendSMS(ctx, wenova.SendSMSRequest{
		Header:      "WNV-OTP",
		PhoneNumber: "2012345678",
		Message:     "Code: 123456",
		Token:       "your-token",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(smsResult)
}
```

## Packages

The package currently supports:

- `SendSMS`

## Configuration

Set `Wenova.BaseUrl` to customize the API host. If it is empty, the default base URL is `https://apimicroservices.wenova.fun`.

## Structure

Like `envtools`, `wenova` uses a flat package structure:

- `wenova.go` defines the package type
- `sms.go` contains SMS helpers
- `wenova_test.go` covers shared behavior

You can call helpers either from package-level functions or from a `Wenova` value.

## SMS

`SendSMS` sends SMS requests to `POST /sms/package`.

### Request Rules

`SendSMS` requires `Token`. If it is missing, it returns an error.

### Example

```go
result, err := wenova.SendSMS(ctx, wenova.SendSMSRequest{
	Header:      "WNV-OTP",
	PhoneNumber: "2012345678",
	Message:     "Code: 123456",
	Token:       "your-token",
})
```

## Error Behavior

- return decoded JSON as `any` on success
- return `nil, nil` for empty successful response bodies
- parse API error payloads and include the HTTP status code in returned errors

## Testing

Run all tests with:

```bash
go test ./...
```

The live Wenova SMS test reads:

- `WENOVA_TOKEN`
- `WENOVA_DEMO_PHONE_NUMBER`
- `WENOVA_DEMO_HEADER` (optional)
- `WENOVA_DEMO_MESSAGE` (optional)

The SMS test requires `WENOVA_TOKEN` and `WENOVA_DEMO_PHONE_NUMBER`. In GitHub Actions, provide these values as repository secrets with the same names.
