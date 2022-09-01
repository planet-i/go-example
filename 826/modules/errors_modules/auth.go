package errors_modules

import "github.com/pkg/errors"

var (
	AuthError       = errors.New("auth failed")
	PermissionError = errors.New("permission error")
)
