package registration

import (
	"context"
	"net/netip"
	"os"
	"time"

	"gitlab.mreg.io/my-registry/auth/domain/registration"
	"gitlab.mreg.io/my-registry/auth/domain/session"
)

type Service interface {
	CreateRegistrationFlow(ctx context.Context, ipAddress netip.Addr, userAgent string) (*registration.Flow, *session.Session, error)
}

type service struct {
	session          session.Repository
	registrationFlow registration.Repository
	interval         time.Duration
}

func NewService(session session.Repository, registrationFlow registration.Repository) Service {
	interval, err := time.ParseDuration(os.Getenv("EXPIRE_INTERVAL"))
	if err != nil {
		panic("Environmental variable EXPIRE_INTERVAL could not be parsed")
	}
	return &service{session, registrationFlow, interval}
}

func (s *service) CreateRegistrationFlow(ctx context.Context, ipAddress netip.Addr, userAgent string) (*registration.Flow, *session.Session, error) {
	sessionModel := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 0,
		ExpiryInterval:              s.interval,
		Devices: []session.Device{
			{IPAddress: ipAddress, UserAgent: userAgent, GeoLocation: "(unimplemented)"}, // TODO ip2Geolocation
		},
	}
	if err := s.session.CreateSession(ctx, sessionModel); err != nil {
		return nil, nil, err
	}
	flow := &registration.Flow{SessionID: sessionModel.ID}
	if err := s.registrationFlow.CreateFlow(ctx, flow); err != nil {
		return nil, nil, err
	}
	return flow, sessionModel, nil
}
