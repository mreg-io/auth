package registration

import "errors"

var (
	ErrEmailExists      = errors.New("email already exists")
	ErrInsecurePassword = errors.New("insecure password")
	ErrSessionExpired   = errors.New("session expired")
	ErrFlowExpired      = errors.New("flow expired")
)
