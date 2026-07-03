package tools

import "fmt"

// validationErr is a parameter validation failure (SPEC 03: "Si falla
// validación -> error 400, no cachear"). It is never sent to the API or cache.
type validationErr struct{ msg string }

func (e *validationErr) Error() string { return e.msg }

func validationError(format string, args ...any) error {
	return &validationErr{msg: fmt.Sprintf(format, args...)}
}
