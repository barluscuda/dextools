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
	client := wenova.Wenova{}

	otpResult, err := client.SendOtp(ctx, wenova.SendOtpRequest{
		Header:      "WNV-OTP",
		PhoneNumber: "2012345678",
		Message:     "Code: 123456",
		Token:       "your-token",
		UsePackage:  true,
	})
	if err != nil {
		panic(err)
	}

	provinces, err := client.GetProvinces(ctx, wenova.Options{
		PluginKey: "your-plugin-key",
		Lang:      "en",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(otpResult)
	fmt.Println(provinces)
}
```

## Packages

The package currently supports:

- `SendOtp`
- `ScriptID`
- `GetProvinces`
- `GetProvinceById`
- `GetDistrictsByProvince`
- `GetDistrictById`
- `GetVillagesByDistrict`
- `GetVillageById`

## Configuration

All HTTP calls use the same base URL resolution rules:

- an explicit `BaseURL` field wins when provided
- otherwise `WENOVA_API_URL` is used when set
- otherwise the default base URL is `https://apimicroservices.wenova.fun`

Example environment value:

```text
WENOVA_API_URL=https://apimicroservices.wenova.fun
```

## Structure

Like `envtools`, `wenova` uses a flat package structure:

- `wenova.go` defines the package type
- `apiclient.go` holds shared base URL and error helpers
- `smsotp.go` contains SMS and OTP helpers
- `address.go` contains address lookup helpers
- `wenova_test.go` covers shared behavior

You can call helpers either from package-level functions or from a `Wenova` value.

## SMS OTP

`SendOtp` sends OTP and SMS package requests to `POST /sms/package`.

### Request Rules

`SendOtp` requires at least one of:

- `Token`
- `ScriptID`

If both are missing, it returns an error.

### Example

```go
result, err := wenova.SendOtp(ctx, wenova.SendOtpRequest{
	Header:      "WNV-OTP",
	PhoneNumber: "2012345678",
	Message:     "Code: 123456",
	Token:       "your-token",
	UsePackage:  true,
})
```

### Helpers

The package also includes:

- `ScriptID(string) int64` for parsing a positive script ID from a string

## Address

The address helpers fetch Wenova Link location data for provinces, districts, and villages.

### Options

`Options` supports:

- `PluginKey` for API access
- `KW` for keyword filtering
- `Lang` for response language
- `BaseURL` for overriding the API host

### Available Functions

- `GetProvinces`
- `GetProvinceById`
- `GetDistrictsByProvince`
- `GetDistrictById`
- `GetVillagesByDistrict`
- `GetVillageById`

### Example

```go
districts, err := wenova.GetDistrictsByProvince(ctx, 1, wenova.Options{
	PluginKey: "your-plugin-key",
	KW:        "chan",
	Lang:      "en",
})
```

### Validation

The address helpers enforce:

- `PluginKey` must not be empty
- numeric IDs must be positive

## Error Behavior

Both request areas:

- return decoded JSON as `any` on success
- return `nil, nil` for empty successful response bodies
- parse API error payloads and include the HTTP status code in returned errors

## Testing

Run all tests with:

```bash
go test ./...
```

The live Wenova address test reads `WENOVA_PLUGIN_KEY`. If it is not set, the test is skipped.

The live Wenova service API test reads:

- `WENOVA_TOKEN`
- `WENOVA_DEMO_PHONE_NUMBER`
- `WENOVA_DEMO_HEADER` (optional)
- `WENOVA_DEMO_MESSAGE` (optional)

If `WENOVA_TOKEN` or `WENOVA_DEMO_PHONE_NUMBER` is missing, the service test is skipped. In GitHub Actions, provide these values as repository secrets with the same names.
