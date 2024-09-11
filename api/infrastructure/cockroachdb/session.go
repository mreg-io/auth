package cockroachdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.mreg.io/my-registry/auth/domain/session"
)

type sessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) session.Repository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session session.Session) (session.Session, error) {
	// TODO implement me
	panic("implement me")
}
