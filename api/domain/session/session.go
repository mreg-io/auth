package session

import (
	"time"
)

type Session struct {
	ID                          string
	Active                      bool
	AuthenticatorAssuranceLevel uint8
	IssuedAt                    time.Time
	ExpiresAt                   time.Time
	AuthenticatedAt             time.Time
	Devices                     []Device

	ExpiryInterval time.Duration

	csrfToken string
}
