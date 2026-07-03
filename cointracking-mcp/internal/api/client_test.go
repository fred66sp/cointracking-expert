package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"testing"
)

func TestSignRequestBodyMatchesStdlibHMAC(t *testing.T) {
	body := "method=getBalance&nonce=1"
	secret := "secret"

	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(body))
	want := hex.EncodeToString(mac.Sum(nil))

	got := SignRequestBody(body, secret)
	if got != want {
		t.Fatalf("SignRequestBody mismatch:\n got  %s\n want %s", got, want)
	}
	if len(got) != 128 {
		t.Fatalf("expected 128 hex chars (SHA-512), got %d", len(got))
	}
}

func TestSignRequestBodyDiffersByInput(t *testing.T) {
	a := SignRequestBody("body-a", "secret")
	b := SignRequestBody("body-b", "secret")
	if a == b {
		t.Fatal("different bodies must not produce the same signature")
	}
	c := SignRequestBody("body-a", "other-secret")
	if a == c {
		t.Fatal("different secrets must not produce the same signature")
	}
}

func TestNextNonceMonotonic(t *testing.T) {
	c := NewClient("key", "secret", nil)
	var last int64 = -1
	for i := 0; i < 1000; i++ {
		n := c.nextNonce()
		parsed, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			t.Fatalf("nonce not numeric: %s", n)
		}
		if parsed <= last {
			t.Fatalf("nonce not strictly increasing: prev=%d got=%d", last, parsed)
		}
		last = parsed
	}
}
