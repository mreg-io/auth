package session

import (
	"net/netip"

	"github.com/jackc/pgx/v5/pgtype"
)

type Device struct {
	ID          pgtype.UUID
	IPAddress   netip.Addr
	GeoLocation string
	UserAgent   string
	SessionID   pgtype.UUID
}
