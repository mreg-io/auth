package identity

import "time"

type Email struct {
	Value      string
	Verified   bool
	VerifiedAt time.Time
	CreateTime time.Time
	UpdateTime time.Time
}

func (e *Email) ETag() (string, error) {
	return "", nil // Placeholder for actual implementation
}
