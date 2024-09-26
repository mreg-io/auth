package identity

import "time"

type IDState int32

const (
	StateActive IDState = iota + 1
	StateSuspended
)

type Identity struct {
	ID              string
	State           IDState
	FullName        string
	DisplayName     string
	AvatarURL       string
	Emails          []Email
	Timezone        string
	CreateTime      time.Time
	UpdateTime      time.Time
	StateUpdateTime time.Time
}

func (i *Identity) ETag() (string, error) {
	return "", nil // Placeholder for actual implementation. For now, ETag is not used in this example.
}
