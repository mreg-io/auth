package registration

import (
	"context"
	"net/netip"
	"os"
	"time"

	"gitlab.mreg.io/my-registry/auth/domain/identity"

	"gitlab.mreg.io/my-registry/auth/domain/registration"
	"gitlab.mreg.io/my-registry/auth/domain/session"
)

type Service interface {
	CreateRegistrationFlow(ctx context.Context, ipAddress netip.Addr, userAgent string) (*registration.Flow, *session.Session, error)
	CompleteRegistrationFlow(context.Context, *registration.Flow, netip.Addr, string) (*session.Session, error)
}

type service struct {
	session              session.Repository
	registrationFlow     registration.Repository
	identityRepo         identity.Repository
	sessionInterval      time.Duration
	registrationInterval time.Duration
}

func NewService(session session.Repository, registrationFlow registration.Repository, identityRepo identity.Repository) Service {
	sessionInterval, err := time.ParseDuration(os.Getenv("SESSION_EXPIRY_INTERVAL"))
	if err != nil {
		panic("Environmental variable SESSION_EXPIRY_INTERVAL could not be parsed")
	}
	registrationInterval, err := time.ParseDuration(os.Getenv("REGISTRATION_EXPIRY_INTERVAL"))
	if err != nil {
		panic("Environmental variable REGISTRATION_EXPIRY_INTERVAL could not be parsed")
	}
	return &service{session, registrationFlow, identityRepo, sessionInterval, registrationInterval}
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

// CompleteRegistrationFlow will fill the flow identity in the process
func (s *service) CompleteRegistrationFlow(ctx context.Context, flow *registration.Flow, ipAddress netip.Addr, userAgent string) (*session.Session, error) {
	var err error

	// check if flow expires
	if err = s.registrationFlow.QueryFlow(ctx, flow); err != nil {
		return nil, err
	}
	if flow.IsExpired() {
		return nil, ErrFlowExpired
	}

	// check if session expired
	preSessionData := &session.Session{ID: flow.SessionID}
	if err = s.session.QuerySession(ctx, preSessionData); err != nil {
		return nil, err
	}
	if preSessionData.IsExpired() {
		return nil, ErrSessionExpired
	}

	UserDevice := &session.Device{IPAddress: ipAddress, UserAgent: userAgent, GeoLocation: "(unimplemented)"}
	if !preSessionData.DeviceExists(UserDevice) {
		if err = s.session.UpdateDevice(ctx, preSessionData, UserDevice); err != nil {
			return nil, err
		}
	}

	// check if email exists
	exist, err := s.identityRepo.EmailExists(ctx, flow.Identity.Emails[0].Value)
	// look up the email in the database
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrEmailExists
	}

	// check if password is insecure
	if !identity.IsSecure(flow.Password) {
		return nil, ErrInsecurePassword
	}

	// create identity
	newIdentity := &identity.Identity{
		State:    identity.StateActive,
		Emails:   []identity.Email{flow.Identity.Emails[0]},
		Timezone: flow.Identity.Timezone,
	}
	// password hash
	newIdentity.PasswordHash, err = identity.CreateHash(flow.Password, identity.DefaultParams)
	if err != nil {
		return nil, err
	}

	if err = s.identityRepo.CreateIdentity(ctx, newIdentity); err != nil {
		return nil, err
	}

	// create session
	sessionModel := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 1,
		ExpiryInterval:              s.sessionInterval,
		Devices: []session.Device{
			{IPAddress: ipAddress, UserAgent: userAgent, GeoLocation: "(unimplemented)"}, // TODO ip2Geolocation
		},
		Identity: newIdentity,
	}
	if err := s.session.CreateSession(ctx, sessionModel); err != nil {
		return nil, err
	}

	return sessionModel, nil
}
