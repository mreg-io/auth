package registration

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Flow struct {
	FlowID    pgtype.UUID
	IssuedAt  pgtype.Timestamptz
	ExpiresAt pgtype.Timestamptz
	SessionID pgtype.UUID
}
