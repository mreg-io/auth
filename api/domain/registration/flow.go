package registration

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"strconv"
	"time"

	"gitlab.mreg.io/my-registry/auth/domain/identity"
)

type Flow struct {
	FlowID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
	SessionID string
	Password  string
	Interval  time.Duration
	Identity  *identity.Identity
}

var crcTable = crc32.MakeTable(crc32.IEEE)

func (f *Flow) ETag() (string, error) {
	if f.SessionID == "" {
		return "", fmt.Errorf("SessionID cannot be empty")
	}
	if f.ExpiresAt.IsZero() {
		return "", fmt.Errorf("ExpiresAt cannot be zero")
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(f.ExpiresAt); err != nil {
		fmt.Println("Error encoding ExpiresAt:", err)
		return "", err
	}
	if err := encoder.Encode(f.SessionID); err != nil {
		fmt.Println("Error encoding SessionID:", err)
		return "", err
	}
	// Compute the CRC32 checksum
	checksum := crc32.Checksum(buffer.Bytes(), crcTable)

	// Convert the checksum to a hexadecimal string
	return fmt.Sprintf("W/\"%s\"", strconv.FormatUint(uint64(checksum), 16)), nil
}

func (f *Flow) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}
