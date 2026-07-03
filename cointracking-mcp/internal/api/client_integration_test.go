package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func withFakeAPI(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	prev := apiURL
	apiURL = srv.URL
	t.Cleanup(func() { apiURL = prev })

	return NewClient("test-key", "test-secret", NewRateTracker(20))
}

func TestClientRequestSuccessVerifiesAuthHeaders(t *testing.T) {
	var gotKey, gotSign, gotMethod string
	c := withFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("Key")
		gotSign = r.Header.Get("Sign")
		_ = r.ParseForm()
		gotMethod = r.Form.Get("method")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":1,"balance":{"BTC":"1.5"}}`))
	})

	var out map[string]any
	if err := c.Request("getBalance", nil, &out); err != nil {
		t.Fatalf("Request: %v", err)
	}
	if gotKey != "test-key" {
		t.Errorf("expected Key header 'test-key', got %q", gotKey)
	}
	if gotSign == "" || len(gotSign) != 128 {
		t.Errorf("expected 128-char hex Sign header, got %q", gotSign)
	}
	if gotMethod != "getBalance" {
		t.Errorf("expected method=getBalance in body, got %q", gotMethod)
	}
}

func TestClientRequestRateLimit(t *testing.T) {
	c := withFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	err := c.Request("getTrades", nil, nil)
	if err == nil {
		t.Fatal("expected rate limit error")
	}
	ctErr, ok := err.(*CTError)
	if !ok || ctErr.Code != CodeRateLimit {
		t.Fatalf("expected CTError with code RATE_LIMIT, got %v", err)
	}
}

func TestClientRequestAPIErrorEnvelope(t *testing.T) {
	c := withFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":0,"error":"Invalid API key"}`))
	})
	err := c.Request("getGains", nil, nil)
	if err == nil {
		t.Fatal("expected API-level error")
	}
	ctErr, ok := err.(*CTError)
	if !ok || ctErr.Code != CodeAPIError {
		t.Fatalf("expected CTError with code API_ERROR, got %v", err)
	}
}

func TestClientRequestErrorsAreNotCacheable(t *testing.T) {
	// Sanity check that RequestRaw returns nil raw on error, so callers
	// (cachedCall) never accidentally cache an error payload.
	c := withFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	})
	raw, err := c.RequestRaw("getBalance", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if raw != nil {
		t.Fatalf("expected nil raw payload on error, got %s", raw)
	}
	var ctErr *CTError
	if !isCTError(err, &ctErr) || ctErr.Code != "HTTP_500" {
		t.Fatalf("expected CTError HTTP_500, got %v", err)
	}
}

func isCTError(err error, out **CTError) bool {
	ce, ok := err.(*CTError)
	if ok {
		*out = ce
	}
	return ok
}
