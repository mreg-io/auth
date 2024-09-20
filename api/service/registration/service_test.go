package registration

import (
	"context"
	"net/netip"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.mreg.io/my-registry/auth/domain/registration"
	"gitlab.mreg.io/my-registry/auth/domain/session"
)

type mockSessionRepository struct {
	mock.Mock
}

func (m *mockSessionRepository) CreateSession(ctx context.Context, session *session.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

type mockFlowRepository struct {
	mock.Mock
}

func (m *mockFlowRepository) CreateFlow(ctx context.Context, flow *registration.Flow) error {
	args := m.Called(ctx, flow)
	return args.Error(0)
}

type serviceTestSuite struct {
	suite.Suite
	service               Service
	mockSessionRepository *mockSessionRepository
	mockFlowRepository    *mockFlowRepository
}

func (s *serviceTestSuite) SetupSuite() {
	s.T().Setenv("SESSION_EXPIRY_INTERVAL", "1h")
	s.T().Setenv("REGISTRATION_EXPIRY_INTERVAL", "1h")
	s.T().Setenv("CSRF_SECRET", "Wryyyyyyyy")
	s.mockSessionRepository = new(mockSessionRepository)
	s.mockFlowRepository = new(mockFlowRepository)
	s.service = NewService(s.mockSessionRepository, s.mockFlowRepository)
}

func (s *serviceTestSuite) TestCreateRegistrationFlow() {
	IPAddress, _ := netip.ParseAddr("192.168.1.1")
	userAgent := "Tor"
	sessionID := "123456789"
	interval, err := time.ParseDuration(os.Getenv("SESSION_EXPIRY_INTERVAL"))
	s.Require().NoError(err)
	issuedAt, err := time.Parse(time.UnixDate, "Wed Feb 25 11:06:39 PST 1069")
	s.Require().NoError(err)
	expiresAt, err := time.Parse(time.UnixDate, "Wed Feb 26 11:06:39 PST 2069")
	s.Require().NoError(err)
	sessionModel := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 0,
		ExpiryInterval:              interval,
		Devices: []session.Device{
			{IPAddress: IPAddress, UserAgent: userAgent, GeoLocation: "(unimplemented)"}, // TODO ip2Geolocation
		},
	}
	registrationInterval, err := time.ParseDuration(os.Getenv("REGISTRATION_EXPIRY_INTERVAL"))
	s.Require().NoError(err)
	flow := &registration.Flow{SessionID: sessionID, Interval: registrationInterval}
	ctx := context.Background()
	call1 := s.mockSessionRepository.
		On("CreateSession", ctx, sessionModel).
		Run(func(args mock.Arguments) {
			sessionModel := args.Get(1).(*session.Session)
			sessionModel.ID = "123456789"
			sessionModel.IssuedAt = issuedAt
			sessionModel.ExpiresAt = expiresAt
		}).
		Return(nil).
		Once()

	call2 := s.mockFlowRepository.
		On("CreateFlow", ctx, flow).
		Run(func(args mock.Arguments) {
			flow := args.Get(1).(*registration.Flow) // Extract the flow from arguments
			flow.FlowID = "987654321"                // Set FlowID
			flow.IssuedAt = issuedAt                 // Set IssuedAt
		}).
		Return(nil).
		Once()
	flow, sessionModel, err = s.service.CreateRegistrationFlow(ctx, IPAddress, userAgent)
	s.Require().NoError(err)
	// Assert proper call to repository
	s.mockSessionRepository.AssertExpectations(s.T())
	s.NotEmpty(flow)
	s.NotEmpty(sessionModel.ID)
	s.NotEmpty(sessionModel.ExpiresAt)
	s.NotEmpty(sessionModel.GetCSRFToken())
	// Assert flow, you should add more
	csrfContent, err := sessionModel.GetCSRFToken()
	s.Require().NoError(err)
	parts := strings.Split(csrfContent, ".")
	s.True(session.VerifyCSRFToken(parts[1], parts[0]))

	// Reset mock
	call1.Unset()
	call2.Unset()
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}
