package registration

import (
	"gitlab.mreg.io/my-registry/auth/session"
)

type service struct {
	repository     Repository
	sessionService session.Service
}

type Service interface {
	CreateRegistrationFlow() (*Flow, error)
}

func NewService(repository Repository, sessionService session.Service) Service {
	return &service{repository, sessionService}
}

func (s *service) CreateRegistrationFlow() (*Flow, error) {
	// Dependency injection
	//s.sessionService.CreateSession()
	flow := &Flow{}
	sessionFlow, err := s.sessionService.CreateSession()
	if err != nil {
		return nil, err
	}
	flow.sessionId = sessionFlow.GetID()

	err = s.repository.insertFlow(flow)
	return flow, err
}
