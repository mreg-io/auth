package session

import (
	"context"
)

type Repository interface {
	CreateSession(ctx context.Context, session Session) (Session, error)
}
