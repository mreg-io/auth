package session

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/fxamacker/cbor/v2"

	"gitlab.mreg.io/my-registry/auth/domain/identity"
)

type Session struct {
	ID                          string `cbor:"1, keyasint"`
	Active                      bool   `cbor:"2, keyasint"`
	AuthenticatorAssuranceLevel uint8  `cbor:"3, keyasint"`
	IssuedAt                    time.Time
	ExpiresAt                   time.Time `cbor:"4, keyasint"`
	AuthenticatedAt             time.Time
	Devices                     []Device `cbor:"5, keyasint,toarray"`
	Identity                    *identity.Identity

	ExpiryInterval time.Duration
}

func (s *Session) ETag() (string, error) {
	if s.ID == "" {
		return "", fmt.Errorf("session ID cannot be empty")
	}

	if s.IssuedAt.IsZero() {
		return "", fmt.Errorf("session must have a IssuedAt time")
	}

	if s.ExpiresAt.IsZero() {
		return "", fmt.Errorf("session must have a ExpiresAt time")
	}
	if len(s.Devices) == 0 {
		return "", fmt.Errorf("session must have a Device")
	}

	// Serialize fields to CBOR
	var buffer bytes.Buffer
	if err := cbor.NewEncoder(&buffer).Encode(s); err != nil {
		return "", err
	}

	// Create CRC32 checksum using IEEE CRC32 table
	checksum := crc32.Checksum(buffer.Bytes(), crc32.MakeTable(crc32.IEEE))
	// Return the checksum as a string
	return fmt.Sprintf("W/\"%x\"", checksum), nil
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
