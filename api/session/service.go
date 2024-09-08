package session

type service struct {
	session    Session
	repository Repository
}

type Service interface {
	CreateSession() (*Session, error)
}

func NewService() Service {
	return &service{}
}
func (s *service) CreateSession() (*Session, error) {
	//TODO
	return &Session{}, nil
}
