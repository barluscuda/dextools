package wenowa

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const DefaultBaseURL = "https://apimicroservices.wenova.fun"

// ResolveBaseURL returns override if set, else WENOVA_API_URL, else DefaultBaseURL.
func ResolveBaseURL(override string) string {
	if strings.TrimSpace(override) != "" {
		return strings.TrimSuffix(strings.TrimSpace(override), "/")
	}
	if v := strings.TrimSpace(os.Getenv("WENOVA_API_URL")); v != "" {
		return strings.TrimSuffix(v, "/")
	}
	return strings.TrimSuffix(DefaultBaseURL, "/")
}

// ResolveBaseURL returns override if set, else WENOVA_API_URL, else DefaultBaseURL.
func (Wenowa) ResolveBaseURL(override string) string {
	return ResolveBaseURL(override)
}

// ErrorFromResponseBody parses JSON error bodies like the Node SDKs.
func ErrorFromResponseBody(b []byte, status int, fallback string) error {
	var wrap struct {
		Message any `json:"message"`
		Error   any `json:"error"`
	}
	msg := fallback
	if json.Unmarshal(b, &wrap) == nil {
		if s, ok := stringifyErrPart(wrap.Message); ok && s != "" {
			msg = s
		} else if s, ok := stringifyErrPart(wrap.Error); ok && s != "" {
			msg = s
		}
	}
	return fmt.Errorf("%s (HTTP %d)", msg, status)
}

func stringifyErrPart(v any) (string, bool) {
	switch x := v.(type) {
	case string:
		return x, true
	case map[string]any:
		if m, ok := x["message"].(string); ok {
			return m, true
		}
		if raw, err := json.Marshal(x); err == nil {
			return string(raw), true
		}
	case float64, bool, json.Number:
		return fmt.Sprint(x), true
	}
	return "", false
}
