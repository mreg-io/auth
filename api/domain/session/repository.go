package session

import (
	"context"
)

type Repository interface {
	CreateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, sessionID string) error
	QuerySession(ctx context.Context, session *Session) error
	UpdateDevice(ctx context.Context, session *Session, newDevice *Device) error
}
