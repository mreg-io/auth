package session

import (
	"time"

	"gitlab.mreg.io/my-registry/auth/domain/identity"
)

type Session struct {
	ID                          string
	Active                      bool
	AuthenticatorAssuranceLevel uint8
	IssuedAt                    time.Time
	ExpiresAt                   time.Time
	AuthenticatedAt             time.Time
	Devices                     []Device
	Identity                    *identity.Identity

	ExpiryInterval time.Duration
}

func (s *Session) ETag() (string, error) {
	panic("not implemented")
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) DeviceExists(userDevice *Device) bool {
	for _, device := range s.Devices {
		if device.IPAddress == userDevice.IPAddress && device.UserAgent == userDevice.UserAgent {
			return true
		}
	}
	return false
}
