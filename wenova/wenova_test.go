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

func TestNewWenovaAPI(t *testing.T) {
	client := NewWenovaAPI(" demo-token ")
	if client.BaseUrl != DefaultBaseURL {
		t.Fatalf("expected default base url %q, got %q", DefaultBaseURL, client.BaseUrl)
	}
	if client.Token != "demo-token" {
		t.Fatalf("expected token on client, got %q", client.Token)
	}
}

func TestSysChangeBaseUrl(t *testing.T) {
	client := NewWenovaAPI("demo-token")
	client.SysChangeBaseUrl(" https://example.com/ ")
	if client.BaseUrl != "https://example.com" {
		t.Fatalf("expected trimmed base url, got %q", client.BaseUrl)
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

func TestSendSMSLive(t *testing.T) {
	phoneNumber := requiredTrimmedEnv("WENOVA_DEMO_PHONE_NUMBER")
	if phoneNumber == "" {
		t.Skip("skipping live Wenova SMS test: WENOVA_DEMO_PHONE_NUMBER is not set")
	}

	token := requiredTrimmedEnv("WENOVA_TOKEN")
	if token == "" {
		t.Skip("skipping live Wenova SMS test: WENOVA_TOKEN is not set")
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

	client := NewWenovaAPI(token)
	result, err := client.SendSMS(ctx, SendSMSRequest{
		Header:      header,
		PhoneNumber: phoneNumber,
		Message:     message,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected SMS response, got nil")
	}
}
