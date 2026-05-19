package wenova

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// SendSMSRequest matches the Wenova SMS payload.
type SendSMSRequest struct {
	Header      string
	PhoneNumber string
	Message     string
	Token       string
}

// SendSMS POSTs to /sms/package and returns the decoded JSON body.
func SendSMS(ctx context.Context, req SendSMSRequest) (any, error) {
	return Wenova{}.SendSMS(ctx, req)
}

func (w Wenova) SendSMS(ctx context.Context, req SendSMSRequest) (any, error) {
	if strings.TrimSpace(req.Token) == "" {
		return nil, errors.New("token is required")
	}

	base := strings.TrimSuffix(strings.TrimSpace(w.BaseUrl), "/")
	if base == "" {
		base = DefaultBaseURL
	}
	url := base + "/sms/package"

	body := map[string]any{
		"header":      req.Header,
		"phoneNumber": req.PhoneNumber,
		"message":     req.Message,
		"token":       strings.TrimSpace(req.Token),
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	hreq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, ErrorFromResponseBody(b, resp.StatusCode, "Failed to send SMS")
	}

	var out any
	if len(b) == 0 {
		return nil, nil
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return out, nil
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
