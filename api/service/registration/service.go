package registration

import (
	"context"
	"net"

	"gitlab.mreg.io/my-registry/auth/domain/registration"
	"gitlab.mreg.io/my-registry/auth/domain/session"
)

type Service interface {
	CreateRegistrationFlow(ctx context.Context, ipAddress net.Addr, userAgent string) (registration.Flow, session.Session, error)
}

type service struct {
	sessionRepository session.Repository
}

func NewService(sessionRepository session.Repository) Service {
	return &service{}
}

func (s *service) CreateRegistrationFlow(ctx context.Context, ipAddress net.Addr, userAgent string) (registration.Flow, session.Session, error) {
	// TODO implement me
	panic("implement me")
}
