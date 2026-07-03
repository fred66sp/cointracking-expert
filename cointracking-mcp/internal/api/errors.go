package api

import "fmt"

// Error codes, mirroring reference-js-repo.md's utils/errors.ts.
const (
	CodeNetwork            = "NETWORK"
	CodeRateLimit          = "RATE_LIMIT"
	CodeBadJSON            = "BAD_JSON"
	CodeAPIError           = "API_ERROR"
	CodeMissingCredentials = "MISSING_CREDENTIALS"
)

// CTError is a CoinTracking API error with an actionable message and a
// machine-readable code for callers to branch on (e.g. rate limit handling).
type CTError struct {
	Message string
	Code    string
}

func (e *CTError) Error() string {
	return fmt.Sprintf("CoinTracking API error [%s]: %s", e.Code, e.Message)
}

func newCTError(code, format string, args ...any) *CTError {
	return &CTError{Code: code, Message: fmt.Sprintf(format, args...)}
}
