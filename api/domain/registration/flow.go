package registration

import (
	"time"
)

type Flow struct {
	FlowID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
	SessionID string
	Interval  time.Duration
}
