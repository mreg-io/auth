package session

type service struct {
	repository Repository
}

type Service interface {
	CreateSession() (*Session, error)
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) CreateSession() (*Session, error) {
	// TODO
	return &Session{}, nil
}
