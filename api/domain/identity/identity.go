package identity

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/fxamacker/cbor/v2"
)

type IDState int32

const (
	StateActive IDState = iota + 1
	StateSuspended
)

type Identity struct {
	ID              string  `cbor:"1, keyasint"`
	State           IDState `cbor:"2, keyasint"`
	FullName        string  `cbor:"3, keyasint, omitempty"`
	DisplayName     string  `cbor:"4, keyasint, omitempty"`
	AvatarURL       string  `cbor:"5, keyasint, omitempty"`
	Emails          []Email `cbor:"6, keyasint, toarray"`
	Timezone        string  `cbor:"7, keyasint, omitempty"`
	CreateTime      time.Time
	UpdateTime      time.Time `cbor:"8, keyasint"`
	StateUpdateTime time.Time
	PasswordHash    string
}

func (i *Identity) ETag() (string, error) {
	if i.ID == "" {
		return "", fmt.Errorf("email value cannot be empty")
	}
	if i.CreateTime.IsZero() {
		return "", fmt.Errorf("email must have a create time")
	}
	if len(i.Emails) == 0 {
		return "", fmt.Errorf("identity must have at least one email")
	}

	// Serialize fields to CBOR
	var buffer bytes.Buffer
	if err := cbor.NewEncoder(&buffer).Encode(i); err != nil {
		return "", err
	}

	// Create CRC32 checksum using IEEE CRC32 table
	checksum := crc32.Checksum(buffer.Bytes(), crc32.MakeTable(crc32.IEEE))
	// Return the checksum as a string
	return fmt.Sprintf("W/\"%x\"", checksum), nil
}
