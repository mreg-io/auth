package session

import (
	"context"
)

type Repository interface {
	CreateSession(ctx context.Context, session *Session) error
	QuerySessionByID(ctx context.Context, session *Session) error
	QuerySessionWithDevices(ctx context.Context, session *Session) error
	InsertDevice(ctx context.Context, newDevice *Device) error
}
