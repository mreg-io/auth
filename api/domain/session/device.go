package session

import (
	"net/netip"
)

type Device struct {
	ID          string
	IPAddress   netip.Addr
	GeoLocation string
	UserAgent   string
	SessionID   string
}
