package util

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Verifier interface {
	Verify(ctx context.Context, token, remoteIP string) (*TurnstileResp, error)
}

type TurnstileResp struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Action      string   `json:"action"`
	CData       string   `json:"cdata"`
}

/********** NOOP (bypass) **********/
type NoopVerifier struct{}

func (NoopVerifier) Verify(ctx context.Context, token, remoteIP string) (*TurnstileResp, error) {
	return &TurnstileResp{Success: true}, nil
}

/********** Cloudflare **********/
type TurnstileVerifier struct {
	Secret string
	Client *http.Client
}

const siteVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

func (v TurnstileVerifier) Verify(ctx context.Context, token, remoteIP string) (*TurnstileResp, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}
	form := url.Values{}
	form.Set("secret", v.Secret)
	form.Set("response", token)
	if remoteIP != "" {
		form.Set("remoteip", remoteIP)
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, siteVerifyURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := v.Client
	if c == nil {
		c = &http.Client{Timeout: 4 * time.Second}
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out TurnstileResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if !out.Success {
		return &out, errors.New("turnstile failed")
	}
	return &out, nil
}

/********** Factory **********/
func NewVerifier(secret string, disabled bool) Verifier {
	if disabled {
		return NoopVerifier{}
	}
	return TurnstileVerifier{Secret: secret}
}
