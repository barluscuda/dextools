package wenova

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func requiredTrimmedEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func optionalBoolEnv(key string) bool {
	v := strings.ToLower(requiredTrimmedEnv(key))
	return v == "1" || v == "true" || v == "yes"
}

func TestResolveBaseURL(t *testing.T) {
	t.Setenv("WENOVA_API_URL", "")
	if got := ResolveBaseURL(""); got != DefaultBaseURL {
		t.Fatalf("expected default base url %q, got %q", DefaultBaseURL, got)
	}

	t.Setenv("WENOVA_API_URL", "https://example.com/")
	if got := ResolveBaseURL(""); got != "https://example.com" {
		t.Fatalf("expected env base url to be trimmed, got %q", got)
	}

	if got := ResolveBaseURL(" https://override.test/ "); got != "https://override.test" {
		t.Fatalf("expected override base url to win, got %q", got)
	}
}

func TestScriptID(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{input: "", want: 0},
		{input: "abc", want: 0},
		{input: "-5", want: 0},
		{input: "42", want: 42},
		{input: " 8 ", want: 8},
	}

	for _, tt := range tests {
		if got := ScriptID(tt.input); got != tt.want {
			t.Fatalf("ScriptID(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestErrorFromResponseBody(t *testing.T) {
	err := ErrorFromResponseBody([]byte(`{"message":"bad request"}`), 400, "fallback")
	if err == nil || err.Error() != "bad request (HTTP 400)" {
		t.Fatalf("unexpected error: %v", err)
	}

	err = ErrorFromResponseBody([]byte(`not-json`), 500, "fallback")
	if err == nil || !strings.Contains(err.Error(), "fallback (HTTP 500)") {
		t.Fatalf("unexpected fallback error: %v", err)
	}
}

func TestGetProvincesLive(t *testing.T) {
	pluginKey := requiredTrimmedEnv("WENOVA_PLUGIN_KEY")
	if pluginKey == "" {
		t.Skip("skipping live Wenova address test: WENOVA_PLUGIN_KEY is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := GetProvinces(ctx, Options{
		PluginKey: pluginKey,
		Lang:      "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected province response, got nil")
	}
}

func TestSendOtpLive(t *testing.T) {
	phoneNumber := requiredTrimmedEnv("WENOVA_DEMO_PHONE_NUMBER")
	if phoneNumber == "" {
		t.Skip("skipping live Wenova service test: WENOVA_DEMO_PHONE_NUMBER is not set")
	}

	token := requiredTrimmedEnv("WENOVA_TOKEN")
	scriptID := ScriptID(requiredTrimmedEnv("WENOVA_SCRIPT_ID"))
	if token == "" && scriptID == 0 {
		t.Skip("skipping live Wenova service test: WENOVA_TOKEN or WENOVA_SCRIPT_ID is required")
	}

	header := requiredTrimmedEnv("WENOVA_DEMO_HEADER")
	if header == "" {
		header = "WNV-OTP"
	}

	message := requiredTrimmedEnv("WENOVA_DEMO_MESSAGE")
	if message == "" {
		message = fmt.Sprintf("Code: %d", time.Now().Unix()%1000000)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := SendOtp(ctx, SendOtpRequest{
		Header:      header,
		PhoneNumber: phoneNumber,
		Message:     message,
		Token:       token,
		ScriptID:    scriptID,
		UsePackage:  optionalBoolEnv("WENOVA_DEMO_USE_PACKAGE"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected OTP response, got nil")
	}
}
