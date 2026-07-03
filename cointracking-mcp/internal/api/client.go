// Package api implements the CoinTracking API v1 client: HMAC-SHA512 auth,
// monotonic nonce, and error classification. Ported from
// cointracking-mcp-main/src/api-client.ts (reference-js-repo.md).
package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// apiURL is a var (not const) so tests can point it at an httptest server.
var apiURL = "https://cointracking.info/api/v1/"

// SetAPIURLForTesting overrides the CoinTracking API base URL and returns a
// restore function. Test-only; not for production use.
func SetAPIURLForTesting(url string) (restore func()) {
	prev := apiURL
	apiURL = url
	return func() { apiURL = prev }
}

// Client performs authenticated requests against the CoinTracking API.
type Client struct {
	APIKey    string
	apiSecret string
	http      *http.Client
	rate      *RateTracker

	nonceMu   sync.Mutex
	lastNonce int64
}

func NewClient(apiKey, apiSecret string, rate *RateTracker) *Client {
	return &Client{
		APIKey:    apiKey,
		apiSecret: apiSecret,
		http:      &http.Client{Timeout: 30 * time.Second},
		rate:      rate,
	}
}

// nextNonce returns a strictly-increasing nonce, guaranteeing uniqueness
// even when two calls land in the same millisecond (mirrors api-client.ts).
func (c *Client) nextNonce() string {
	c.nonceMu.Lock()
	defer c.nonceMu.Unlock()
	n := time.Now().UnixMilli()
	if n <= c.lastNonce {
		n = c.lastNonce + 1
	}
	c.lastNonce = n
	return strconv.FormatInt(n, 10)
}

// SignRequestBody computes the HMAC-SHA512 signature of the URL-encoded
// request body, hex-encoded, as required by the CoinTracking `Sign` header.
func SignRequestBody(body, apiSecret string) string {
	mac := hmac.New(sha512.New, []byte(apiSecret))
	mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}

// RequestRaw performs an authenticated POST to the CoinTracking API and
// returns the raw JSON response bytes (post-envelope validation), suitable
// for caching verbatim and for unmarshaling by the caller.
func (c *Client) RequestRaw(method string, params map[string]string) (json.RawMessage, error) {
	form := url.Values{}
	form.Set("method", method)
	form.Set("nonce", c.nextNonce())
	for k, v := range params {
		if v == "" {
			continue
		}
		form.Set(k, v)
	}

	body := form.Encode()
	sign := SignRequestBody(body, c.apiSecret)

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(body))
	if err != nil {
		return nil, newCTError(CodeNetwork, "building request: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Key", c.APIKey)
	req.Header.Set("Sign", sign)

	if c.rate != nil {
		c.rate.RecordCall(method)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, newCTError(CodeNetwork, "contacting CoinTracking: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, newCTError(CodeRateLimit,
			"Rate limit exceeded. Wait before retrying, and prefer targeted queries with start/end filters and a limit.")
	}

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := string(respBody)
		if len(snippet) > 500 {
			snippet = snippet[:500]
		}
		return nil, newCTError(fmt.Sprintf("HTTP_%d", resp.StatusCode), "%s%s", resp.Status, ifNonEmpty(snippet))
	}

	var envelope struct {
		Success *int   `json:"success"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, newCTError(CodeBadJSON, "invalid JSON response: %s", err)
	}
	if envelope.Success != nil && *envelope.Success == 0 {
		msg := envelope.Error
		if msg == "" {
			msg = "Unknown API error"
		}
		return nil, newCTError(CodeAPIError, "%s (method=%s). Check API key permissions and parameters.", msg, method)
	}

	return json.RawMessage(respBody), nil
}

// Request performs an authenticated call and unmarshals the JSON response
// into out.
func (c *Client) Request(method string, params map[string]string, out any) error {
	raw, err := c.RequestRaw(method, params)
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return newCTError(CodeBadJSON, "invalid JSON response: %s", err)
	}
	return nil
}

func ifNonEmpty(s string) string {
	if s == "" {
		return ""
	}
	return ": " + s
}
