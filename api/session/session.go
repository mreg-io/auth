package session

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Session struct {
	ID                          pgtype.UUID
	IssuedAt                    pgtype.Timestamptz
	ExpiresAt                   pgtype.Timestamptz
	Active                      bool
	AuthenticatedAt             pgtype.Timestamptz
	AuthenticatorAssuranceLevel int16
	CSRFToken                   string
}
