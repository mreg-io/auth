package identity

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/fxamacker/cbor/v2"
)

type Email struct {
	Value      string `cbor:"1, keyasint"`
	Verified   bool   `cbor:"2, keyasint, omitempty"`
	VerifiedAt time.Time
	CreateTime time.Time
	UpdateTime time.Time `cbor:"3,keyasint"`
}

func (e *Email) ETag() (string, error) {
	if e.Value == "" {
		return "", fmt.Errorf("email value cannot be empty")
	}
	if e.CreateTime.IsZero() {
		return "", fmt.Errorf("email must have a create time")
	}

	// Serialize fields to CBOR
	var buffer bytes.Buffer
	if err := cbor.NewEncoder(&buffer).Encode(e); err != nil {
		return "", err
	}

	// Create CRC32 checksum using IEEE CRC32 table
	checksum := crc32.Checksum(buffer.Bytes(), crc32.MakeTable(crc32.IEEE))
	// Return the checksum as a string
	return fmt.Sprintf("W/\"%x\"", checksum), nil
}
