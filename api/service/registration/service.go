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
	CompleteRegistrationFlow(context.Context, *registration.Flow) (*session.Session, error)
}

type service struct {
	session              session.Repository
	registrationFlow     registration.Repository
	sessionInterval      time.Duration
	registrationInterval time.Duration
}

func NewService(session session.Repository, registrationFlow registration.Repository) Service {
	sessionInterval, err := time.ParseDuration(os.Getenv("SESSION_EXPIRY_INTERVAL"))
	if err != nil {
		panic("Environmental variable SESSION_EXPIRY_INTERVAL could not be parsed")
	}
	registrationInterval, err := time.ParseDuration(os.Getenv("REGISTRATION_EXPIRY_INTERVAL"))
	if err != nil {
		panic("Environmental variable REGISTRATION_EXPIRY_INTERVAL could not be parsed")
	}
	return &service{session, registrationFlow, sessionInterval, registrationInterval}
}

func (s *service) CreateRegistrationFlow(ctx context.Context, ipAddress netip.Addr, userAgent string) (*registration.Flow, *session.Session, error) {
	sessionModel := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 0,
		ExpiryInterval:              s.sessionInterval,
		Devices: []session.Device{
			{IPAddress: ipAddress, UserAgent: userAgent, GeoLocation: "(unimplemented)"}, // TODO ip2Geolocation
		},
	}
	if err := s.session.CreateSession(ctx, sessionModel); err != nil {
		return nil, nil, err
	}
	flow := &registration.Flow{SessionID: sessionModel.ID, Interval: s.registrationInterval}
	if err := s.registrationFlow.CreateFlow(ctx, flow); err != nil {
		return nil, nil, err
	}
	return flow, sessionModel, nil
}

func (s *service) CompleteRegistrationFlow(context.Context, *registration.Flow) (*session.Session, error) {
	// TODO: Implement me
	panic("Not implemented")
}
