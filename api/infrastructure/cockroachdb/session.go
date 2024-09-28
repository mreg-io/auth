package cockroachdb

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.mreg.io/my-registry/auth/domain/session"
)

//go:embed sql/createSession.sql
var createSessionSQL string

type sessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) session.Repository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *session.Session) error {
	if len(session.Devices) != 1 {
		return errors.New("only one device is allowed when creating a session")
	}

	device := session.Devices[0]

	return r.db.
		QueryRow(
			ctx,
			createSessionSQL,
			session.Active, zeronull.Int2(session.AuthenticatorAssuranceLevel), session.ExpiryInterval, zeronull.Timestamptz(session.AuthenticatedAt),
			session.Identity,
			device.IPAddress, device.GeoLocation, device.UserAgent,
		).
		Scan(&session.ID, &session.IssuedAt, &session.ExpiresAt, &session.Devices[0].ID)
}

func (r *sessionRepository) DeleteSession(context.Context, string) error {
	// TODO implement me
	panic("implement me")
}

func (r *sessionRepository) QuerySession(context.Context, *session.Session) error {
	// TODO implement me
	panic("implement me")
}

func (r *sessionRepository) UpdateDevice(context.Context, *session.Session, *session.Device) error {
	// TODO implement me
	panic("implement me")
}
