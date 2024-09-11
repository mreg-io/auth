package session

import (
	"time"
)

type Session struct {
	ID                          string
	IssuedAt                    time.Time
	ExpiresAt                   time.Time
	Active                      bool
	AuthenticatedAt             time.Time
	AuthenticatorAssuranceLevel int16
	devices                     []Device

	csrfToken string
}
