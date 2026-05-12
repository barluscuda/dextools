package envtools

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// fallbackCheck handles optional fallback logic for Load methods
func fallbackCheck[T any](val T, ok bool, fallback []T) T {
	if ok {
		return val
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	var zero T
	return zero
}

// ensureValue is a helper for MustLoad methods to handle panics
func ensureValue[T any](key string, val T, ok bool, err error) T {
	if !ok || err != nil {
		reason := "missing"
		if err != nil {
			reason = fmt.Sprintf("invalid (%v)", err)
		}
		panic(fmt.Sprintf("envtools: required environment variable %q is %s", key, reason))
	}
	return val
}

// --- STRING ---

func (EnvTools) LoadString(key string, fallback ...string) string {
	val := os.Getenv(key)
	return fallbackCheck(val, val != "", fallback)
}

func (e EnvTools) MustLoadString(key string) string {
	val := os.Getenv(key)
	return ensureValue(key, val, val != "", nil)
}

// --- INT ---

func (EnvTools) LoadInt(key string, fallback ...int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallbackCheck(0, false, fallback)
	}
	i, err := strconv.Atoi(val)
	return fallbackCheck(i, err == nil, fallback)
}

func (e EnvTools) MustLoadInt(key string) int {
	val := os.Getenv(key)
	i, err := strconv.Atoi(val)
	return ensureValue(key, i, val != "", err)
}

// --- INT64 ---

func (EnvTools) LoadInt64(key string, fallback ...int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return fallbackCheck(int64(0), false, fallback)
	}
	i, err := strconv.ParseInt(val, 10, 64)
	return fallbackCheck(i, err == nil, fallback)
}

func (e EnvTools) MustLoadInt64(key string) int64 {
	val := os.Getenv(key)
	i, err := strconv.ParseInt(val, 10, 64)
	return ensureValue(key, i, val != "", err)
}

// --- BOOL ---

func (EnvTools) LoadBool(key string, fallback ...bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallbackCheck(false, false, fallback)
	}
	b, err := strconv.ParseBool(val)
	return fallbackCheck(b, err == nil, fallback)
}

func (e EnvTools) MustLoadBool(key string) bool {
	val := os.Getenv(key)
	b, err := strconv.ParseBool(val)
	return ensureValue(key, b, val != "", err)
}

// --- FLOAT64 ---

func (EnvTools) LoadFloat64(key string, fallback ...float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return fallbackCheck(0.0, false, fallback)
	}
	f, err := strconv.ParseFloat(val, 64)
	return fallbackCheck(f, err == nil, fallback)
}

func (e EnvTools) MustLoadFloat64(key string) float64 {
	val := os.Getenv(key)
	f, err := strconv.ParseFloat(val, 64)
	return ensureValue(key, f, val != "", err)
}

// --- DURATION ---

func (EnvTools) LoadDuration(key string, fallback ...time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallbackCheck(time.Duration(0), false, fallback)
	}
	d, err := time.ParseDuration(val)
	return fallbackCheck(d, err == nil, fallback)
}

func (e EnvTools) MustLoadDuration(key string) time.Duration {
	val := os.Getenv(key)
	d, err := time.ParseDuration(val)
	return ensureValue(key, d, val != "", err)
}

// --- TIME (RFC3339) ---

func (EnvTools) LoadTime(key string, fallback ...time.Time) time.Time {
	val := os.Getenv(key)
	if val == "" {
		return fallbackCheck(time.Time{}, false, fallback)
	}
	t, err := time.Parse(time.RFC3339, val)
	return fallbackCheck(t, err == nil, fallback)
}

func (e EnvTools) MustLoadTime(key string) time.Time {
	val := os.Getenv(key)
	t, err := time.Parse(time.RFC3339, val)
	return ensureValue(key, t, val != "", err)
}

// --- SLICES ---

func (EnvTools) LoadStringSlice(key string, fallback ...[]string) []string {
	val := os.Getenv(key)
	if val == "" {
		return fallbackCheck([]string{}, false, fallback)
	}
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return fallbackCheck(parts, true, fallback)
}

func (e EnvTools) MustLoadStringSlice(key string) []string {
	val := os.Getenv(key)
	res := e.LoadStringSlice(key)
	return ensureValue(key, res, val != "", nil)
}

// --- BYTES (BASE64) ---

func (EnvTools) LoadBase64(key string, fallback ...[]byte) []byte {
	val := os.Getenv(key)
	if val == "" {
		return fallbackCheck([]byte{}, false, fallback)
	}
	data, err := base64.StdEncoding.DecodeString(val)
	return fallbackCheck(data, err == nil, fallback)
}

func (e EnvTools) MustLoadBase64(key string) []byte {
	val := os.Getenv(key)
	data, err := base64.StdEncoding.DecodeString(val)
	return ensureValue(key, data, val != "", err)
}

// --- SPECIALIZED ---

func (e EnvTools) LoadPort(key string, fallback ...int) int {
	p := e.LoadInt(key, fallback...)
	if p < 1 || p > 65535 {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return 0
	}
	return p
}

func (e EnvTools) MustLoadPort(key string) int {
	p := e.MustLoadInt(key)
	if p < 1 || p > 65535 {
		panic(fmt.Sprintf("envtools: port %q out of range: %d", key, p))
	}
	return p
}