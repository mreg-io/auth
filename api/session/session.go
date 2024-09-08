package session

import (
	"github.com/jackc/pgx/v5/pgtype"
	"net/netip"
)

type Session struct {
	id                          pgtype.UUID
	issuedAt                    pgtype.Timestamptz
	expireAt                    pgtype.Timestamptz
	active                      bool
	authenticatedAt             pgtype.Timestamptz
	authenticatorAssuranceLevel int16
	csrfToken                   string
}

func (s *Session) GetID() pgtype.UUID {
	return s.id
}

type Device struct {
	id          pgtype.UUID
	ipAddress   netip.Addr
	geoLocation string
	userAgent   string
	sessionId   pgtype.UUID
}
