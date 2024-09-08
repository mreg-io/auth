package registration

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Flow struct {
	flowId         pgtype.UUID
	issueTime      pgtype.Timestamptz
	expirationTime pgtype.Timestamptz
	sessionId      pgtype.UUID
}
