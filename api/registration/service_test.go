package registration

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"gitlab.mreg.io/my-registry/auth/session"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) insertFlow(flow *Flow) error {
	// TODO implement me
	panic("implement me")
}

type mockSessionService struct {
	mock.Mock
}

func (m *mockSessionService) CreateSession() (*session.Session, error) {
	// TODO implement me
	panic("implement me")
}

type serviceTestSuite struct {
	suite.Suite

	service        Service
	mockRepository *mockRepository
	sessionService session.Service
}

func (s *serviceTestSuite) SetupSuite() {
	s.mockRepository = new(mockRepository)
	s.sessionService = new(mockSessionService)
	s.service = NewService(s.mockRepository, s.sessionService)
}

func (s *serviceTestSuite) TestCreateRegistrationFlow_WithoutErr() {
	flowID := pgtype.UUID{Bytes: [16]byte([]byte("01J75J7AXNYCCCSDASX845RF5W")), Valid: true}
	call := s.mockRepository.
		On("insertFlow", mock.Anything).
		Run(func(args mock.Arguments) {
			flow := args.Get(0).(*Flow)
			flow.FlowID = flowID
		}).
		Return(nil).
		Once()

	flow, err := s.service.CreateRegistrationFlow()

	// Assert proper call to repository
	s.mockRepository.AssertExpectations(s.T())

	// Assert flow, you should add more
	s.NotEmpty(flow)
	s.Equal(flow.FlowID, flowID)

	// Assert err
	s.Require().NoError(err)

	// Reset mock
	call.Unset()
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}
