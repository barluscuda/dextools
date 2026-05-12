package envtools

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestLoadFunctions(t *testing.T) {
	env := EnvTools{}
	now := time.Date(2026, 5, 12, 10, 30, 0, 0, time.UTC)
	encoded := base64.StdEncoding.EncodeToString([]byte("hello"))

	tests := []struct {
		name string
		key  string
		val  string
		want any
		load func(*testing.T) any
	}{
		{
			name: "string",
			key:  "TEST_LOAD_STRING",
			val:  "value",
			want: "value",
			load: func(t *testing.T) any {
				return env.LoadString("TEST_LOAD_STRING")
			},
		},
		{
			name: "int",
			key:  "TEST_LOAD_INT",
			val:  "42",
			want: 42,
			load: func(t *testing.T) any {
				return env.LoadInt("TEST_LOAD_INT")
			},
		},
		{
			name: "int64",
			key:  "TEST_LOAD_INT64",
			val:  "922337203685477580",
			want: int64(922337203685477580),
			load: func(t *testing.T) any {
				return env.LoadInt64("TEST_LOAD_INT64")
			},
		},
		{
			name: "bool",
			key:  "TEST_LOAD_BOOL",
			val:  "true",
			want: true,
			load: func(t *testing.T) any {
				return env.LoadBool("TEST_LOAD_BOOL")
			},
		},
		{
			name: "float64",
			key:  "TEST_LOAD_FLOAT64",
			val:  "3.14",
			want: 3.14,
			load: func(t *testing.T) any {
				return env.LoadFloat64("TEST_LOAD_FLOAT64")
			},
		},
		{
			name: "duration",
			key:  "TEST_LOAD_DURATION",
			val:  "2h45m",
			want: 2*time.Hour + 45*time.Minute,
			load: func(t *testing.T) any {
				return env.LoadDuration("TEST_LOAD_DURATION")
			},
		},
		{
			name: "time",
			key:  "TEST_LOAD_TIME",
			val:  now.Format(time.RFC3339),
			want: now,
			load: func(t *testing.T) any {
				return env.LoadTime("TEST_LOAD_TIME")
			},
		},
		{
			name: "string slice",
			key:  "TEST_LOAD_STRING_SLICE",
			val:  "alpha, beta,gamma ",
			want: []string{"alpha", "beta", "gamma"},
			load: func(t *testing.T) any {
				return env.LoadStringSlice("TEST_LOAD_STRING_SLICE")
			},
		},
		{
			name: "base64",
			key:  "TEST_LOAD_BASE64",
			val:  encoded,
			want: []byte("hello"),
			load: func(t *testing.T) any {
				return env.LoadBase64("TEST_LOAD_BASE64")
			},
		},
		{
			name: "port",
			key:  "TEST_LOAD_PORT",
			val:  "8080",
			want: 8080,
			load: func(t *testing.T) any {
				return env.LoadPort("TEST_LOAD_PORT")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(tc.key, tc.val)

			got := tc.load(t)
			switch want := tc.want.(type) {
			case []string:
				gotSlice, ok := got.([]string)
				if !ok {
					t.Fatalf("got type %T, want []string", got)
				}
				if len(gotSlice) != len(want) {
					t.Fatalf("got %v, want %v", gotSlice, want)
				}
				for i := range want {
					if gotSlice[i] != want[i] {
						t.Fatalf("got %v, want %v", gotSlice, want)
					}
				}
			case []byte:
				gotBytes, ok := got.([]byte)
				if !ok {
					t.Fatalf("got type %T, want []byte", got)
				}
				if string(gotBytes) != string(want) {
					t.Fatalf("got %q, want %q", gotBytes, want)
				}
			default:
				if got != want {
					t.Fatalf("got %v, want %v", got, want)
				}
			}
		})
	}
}

func TestLoadFunctionsFallback(t *testing.T) {
	env := EnvTools{}
	fallbackTime := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name string
		key  string
		val  string
		want any
		load func() any
	}{
		{
			name: "missing string fallback",
			key:  "TEST_FALLBACK_STRING",
			want: "fallback",
			load: func() any { return env.LoadString("TEST_FALLBACK_STRING", "fallback") },
		},
		{
			name: "invalid int fallback",
			key:  "TEST_FALLBACK_INT",
			val:  "bad",
			want: 7,
			load: func() any { return env.LoadInt("TEST_FALLBACK_INT", 7) },
		},
		{
			name: "invalid bool fallback",
			key:  "TEST_FALLBACK_BOOL",
			val:  "not-bool",
			want: true,
			load: func() any { return env.LoadBool("TEST_FALLBACK_BOOL", true) },
		},
		{
			name: "invalid float fallback",
			key:  "TEST_FALLBACK_FLOAT",
			val:  "oops",
			want: 1.5,
			load: func() any { return env.LoadFloat64("TEST_FALLBACK_FLOAT", 1.5) },
		},
		{
			name: "invalid duration fallback",
			key:  "TEST_FALLBACK_DURATION",
			val:  "later",
			want: 5 * time.Second,
			load: func() any { return env.LoadDuration("TEST_FALLBACK_DURATION", 5*time.Second) },
		},
		{
			name: "invalid time fallback",
			key:  "TEST_FALLBACK_TIME",
			val:  "not-rfc3339",
			want: fallbackTime,
			load: func() any { return env.LoadTime("TEST_FALLBACK_TIME", fallbackTime) },
		},
		{
			name: "invalid base64 fallback",
			key:  "TEST_FALLBACK_BASE64",
			val:  "%%%not-base64%%%",
			want: []byte("fallback"),
			load: func() any { return env.LoadBase64("TEST_FALLBACK_BASE64", []byte("fallback")) },
		},
		{
			name: "out of range port fallback",
			key:  "TEST_FALLBACK_PORT",
			val:  "70000",
			want: 9000,
			load: func() any { return env.LoadPort("TEST_FALLBACK_PORT", 9000) },
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.val != "" {
				t.Setenv(tc.key, tc.val)
			}

			got := tc.load()
			switch want := tc.want.(type) {
			case []byte:
				gotBytes, ok := got.([]byte)
				if !ok {
					t.Fatalf("got type %T, want []byte", got)
				}
				if string(gotBytes) != string(want) {
					t.Fatalf("got %q, want %q", gotBytes, want)
				}
			default:
				if got != want {
					t.Fatalf("got %v, want %v", got, want)
				}
			}
		})
	}
}

func TestMustLoadFunctions(t *testing.T) {
	env := EnvTools{}
	now := time.Date(2026, 5, 12, 10, 30, 0, 0, time.UTC)
	encoded := base64.StdEncoding.EncodeToString([]byte("hello"))

	t.Run("success cases", func(t *testing.T) {
		t.Setenv("TEST_MUST_STRING", "value")
		t.Setenv("TEST_MUST_INT", "42")
		t.Setenv("TEST_MUST_INT64", "99")
		t.Setenv("TEST_MUST_BOOL", "true")
		t.Setenv("TEST_MUST_FLOAT64", "3.14")
		t.Setenv("TEST_MUST_DURATION", "30s")
		t.Setenv("TEST_MUST_TIME", now.Format(time.RFC3339))
		t.Setenv("TEST_MUST_STRING_SLICE", "a, b")
		t.Setenv("TEST_MUST_BASE64", encoded)
		t.Setenv("TEST_MUST_PORT", "8080")

		if got := env.MustLoadString("TEST_MUST_STRING"); got != "value" {
			t.Fatalf("got %q, want %q", got, "value")
		}
		if got := env.MustLoadInt("TEST_MUST_INT"); got != 42 {
			t.Fatalf("got %d, want 42", got)
		}
		if got := env.MustLoadInt64("TEST_MUST_INT64"); got != 99 {
			t.Fatalf("got %d, want 99", got)
		}
		if got := env.MustLoadBool("TEST_MUST_BOOL"); got != true {
			t.Fatalf("got %t, want true", got)
		}
		if got := env.MustLoadFloat64("TEST_MUST_FLOAT64"); got != 3.14 {
			t.Fatalf("got %v, want 3.14", got)
		}
		if got := env.MustLoadDuration("TEST_MUST_DURATION"); got != 30*time.Second {
			t.Fatalf("got %v, want 30s", got)
		}
		if got := env.MustLoadTime("TEST_MUST_TIME"); !got.Equal(now) {
			t.Fatalf("got %v, want %v", got, now)
		}
		if got := env.MustLoadStringSlice("TEST_MUST_STRING_SLICE"); len(got) != 2 || got[0] != "a" || got[1] != "b" {
			t.Fatalf("got %v, want [a b]", got)
		}
		if got := env.MustLoadBase64("TEST_MUST_BASE64"); string(got) != "hello" {
			t.Fatalf("got %q, want %q", got, "hello")
		}
		if got := env.MustLoadPort("TEST_MUST_PORT"); got != 8080 {
			t.Fatalf("got %d, want 8080", got)
		}
	})

	panicTests := []struct {
		name string
		key  string
		val  string
		fn   func()
	}{
		{
			name: "missing string",
			key:  "TEST_PANIC_STRING",
			fn:   func() { env.MustLoadString("TEST_PANIC_STRING") },
		},
		{
			name: "invalid int",
			key:  "TEST_PANIC_INT",
			val:  "bad",
			fn:   func() { env.MustLoadInt("TEST_PANIC_INT") },
		},
		{
			name: "invalid int64",
			key:  "TEST_PANIC_INT64",
			val:  "bad",
			fn:   func() { env.MustLoadInt64("TEST_PANIC_INT64") },
		},
		{
			name: "invalid bool",
			key:  "TEST_PANIC_BOOL",
			val:  "bad",
			fn:   func() { env.MustLoadBool("TEST_PANIC_BOOL") },
		},
		{
			name: "invalid float64",
			key:  "TEST_PANIC_FLOAT64",
			val:  "bad",
			fn:   func() { env.MustLoadFloat64("TEST_PANIC_FLOAT64") },
		},
		{
			name: "invalid duration",
			key:  "TEST_PANIC_DURATION",
			val:  "bad",
			fn:   func() { env.MustLoadDuration("TEST_PANIC_DURATION") },
		},
		{
			name: "invalid time",
			key:  "TEST_PANIC_TIME",
			val:  "bad",
			fn:   func() { env.MustLoadTime("TEST_PANIC_TIME") },
		},
		{
			name: "missing string slice",
			key:  "TEST_PANIC_STRING_SLICE",
			fn:   func() { env.MustLoadStringSlice("TEST_PANIC_STRING_SLICE") },
		},
		{
			name: "invalid base64",
			key:  "TEST_PANIC_BASE64",
			val:  "bad%%%",
			fn:   func() { env.MustLoadBase64("TEST_PANIC_BASE64") },
		},
		{
			name: "out of range port",
			key:  "TEST_PANIC_PORT",
			val:  "70000",
			fn:   func() { env.MustLoadPort("TEST_PANIC_PORT") },
		},
	}

	for _, tc := range panicTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.val != "" {
				t.Setenv(tc.key, tc.val)
			}

			defer func() {
				if recover() == nil {
					t.Fatalf("expected panic")
				}
			}()

			tc.fn()
		})
	}
}
