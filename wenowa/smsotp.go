package wenowa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SendOtpRequest matches the Node SDK payload.
type SendOtpRequest struct {
	Header      string
	PhoneNumber string
	Message     string
	Token       string
	ScriptID    int64
	UsePackage  bool
	BaseURL     string
}

// SendOtp POSTs to /sms/package and returns the decoded JSON body.
func SendOtp(ctx context.Context, req SendOtpRequest) (any, error) {
	return Wenowa{}.SendOtp(ctx, req)
}

func (Wenowa) SendOtp(ctx context.Context, req SendOtpRequest) (any, error) {
	hasToken := strings.TrimSpace(req.Token) != ""
	hasScript := req.ScriptID > 0
	if !hasToken && !hasScript {
		return nil, errors.New("either token or scriptId is required")
	}

	base := ResolveBaseURL(req.BaseURL)
	url := base + "/sms/package"

	body := map[string]any{
		"header":      req.Header,
		"phoneNumber": req.PhoneNumber,
		"message":     req.Message,
		"usePackage":  req.UsePackage,
	}
	if hasToken {
		body["token"] = strings.TrimSpace(req.Token)
	}
	if hasScript {
		body["scriptId"] = req.ScriptID
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
		return nil, ErrorFromResponseBody(b, resp.StatusCode, "Failed to send OTP")
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

// ScriptID parses a positive script id from string.
func ScriptID(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

// ScriptID parses a positive script id from string.
func (Wenowa) ScriptID(s string) int64 {
	return ScriptID(s)
}
