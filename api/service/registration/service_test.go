package registration

import (
	"context"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"gitlab.mreg.io/my-registry/auth/domain/identity"

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

func (m *mockSessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *mockSessionRepository) QuerySessionByID(ctx context.Context, session *session.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *mockSessionRepository) QuerySessionWithDevices(ctx context.Context, session *session.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *mockSessionRepository) InsertDevice(ctx context.Context, newDevice *session.Device) error {
	args := m.Called(ctx, newDevice)
	return args.Error(0)
}

type mockFlowRepository struct {
	mock.Mock
}

func (m *mockFlowRepository) CreateFlow(ctx context.Context, flow *registration.Flow) error {
	args := m.Called(ctx, flow)
	return args.Error(0)
}

func (m *mockFlowRepository) QueryFlowByFlowID(ctx context.Context, flow *registration.Flow) error {
	args := m.Called(ctx, flow)
	return args.Error(0)
}

type mockIdentityRepository struct {
	mock.Mock
}

func (m *mockIdentityRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *mockIdentityRepository) CreateIdentity(ctx context.Context, id *identity.Identity) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockIdentityRepository) QueryEmail(ctx context.Context, email *identity.Email) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

type serviceTestSuite struct {
	suite.Suite
	service                Service
	mockSessionRepository  *mockSessionRepository
	mockFlowRepository     *mockFlowRepository
	mockIdentityRepository *mockIdentityRepository
}

func (s *serviceTestSuite) SetupSuite() {
	s.T().Setenv("SESSION_EXPIRY_INTERVAL", "1h")
	s.T().Setenv("REGISTRATION_EXPIRY_INTERVAL", "1h")
	s.T().Setenv("CSRF_SECRET", "Wryyyyyyyy")
	s.mockSessionRepository = new(mockSessionRepository)
	s.mockFlowRepository = new(mockFlowRepository)
	s.mockIdentityRepository = new(mockIdentityRepository)

	s.service = NewService(s.mockSessionRepository, s.mockFlowRepository, s.mockIdentityRepository)
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

	// Reset mock
	call1.Unset()
	call2.Unset()
}

var (
	userAgent = "Mozilla/5.0"
	sessionID = "123456789"
	email     = "test@example.com"
	password  = "!Securepassword123"
	timezone  = "America/New_York"
)

func (s *serviceTestSuite) TestCompleteRegistrationFlow_Success() {
	ipAddress, _ := netip.ParseAddr("192.168.1.1")
	sessionInterval, err := time.ParseDuration(os.Getenv("SESSION_EXPIRY_INTERVAL"))
	s.Require().NoError(err)
	idCreateTime, err := time.Parse(time.UnixDate, "Wed Feb 25 11:06:39 PST 1069")
	s.Require().NoError(err)
	// Arrange: create a valid flow and session
	registrationFlow := &registration.Flow{
		SessionID: sessionID,
		Identity: &identity.Identity{
			Emails:   []identity.Email{{Value: email}},
			Timezone: timezone,
		},
		Password: password,
	}
	newIdentity := &identity.Identity{
		CreateTime: idCreateTime,
		State:      identity.StateActive,
		Timezone:   timezone,
		Emails:     []identity.Email{{Value: email}}, // too lazy to init other fields
	}
	expectedSession := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 1,
		ExpiryInterval:              sessionInterval,
		Devices: []session.Device{
			{IPAddress: ipAddress, UserAgent: userAgent, GeoLocation: "(unimplemented)"}, // TODO ip2Geolocation
		},
		Identity: newIdentity,
	}
	ctx := context.Background()
	// Mocking expected behavior for valid flow and session
	s.mockFlowRepository.On("QueryFlowByFlowID", ctx, registrationFlow).
		Run(func(args mock.Arguments) {
			registrationFlow := args.Get(1).(*registration.Flow)
			registrationFlow.ExpiresAt = time.Now().Add(900 * time.Hour)
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("QuerySessionWithDevices", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			preSession := args.Get(1).(*session.Session)
			preSession.Devices = []session.Device{
				{GeoLocation: "not the same anyway"},
			}
			preSession.ExpiresAt = time.Now().Add(900 * time.Hour)
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("InsertDevice", ctx, mock.Anything).Return(nil).Once()
	s.mockIdentityRepository.On("EmailExists", ctx, registrationFlow.Identity.Emails[0].Value).Return(false, nil).Once() // Email doesn't exist
	s.mockIdentityRepository.On("CreateIdentity", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			newIdentity := args.Get(1).(*identity.Identity)
			newIdentity.CreateTime = idCreateTime // too lazy to init other fields
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("CreateSession", ctx, mock.Anything).Return(nil).Once()
	// Act: call CompleteRegistrationFlow
	name := "registrationFlows/" + uuid.New().String()
	sessionModel, err := s.service.CompleteRegistrationFlow(ctx, registrationFlow, name, ipAddress, userAgent)
	s.Require().NoError(err)
	expectedSession.Identity = sessionModel.Identity
	// Assert: Ensure no error and valid session returned
	s.Require().NoError(err)
	s.NotNil(sessionModel)
	s.Equal(sessionModel.Identity.Emails[0].Value, email)
	s.Equal(timezone, sessionModel.Identity.Timezone)
	s.Equal(expectedSession, sessionModel)
	// Assert: Ensure mocks were called
	s.mockFlowRepository.AssertExpectations(s.T())
	s.mockSessionRepository.AssertExpectations(s.T())
	s.mockIdentityRepository.AssertExpectations(s.T())
}

func (s *serviceTestSuite) TestCompleteRegistrationFlow_FlowExpire() {
	ipAddress, _ := netip.ParseAddr("192.168.1.1")
	// Arrange: create a valid flow and session
	registrationFlow := &registration.Flow{
		SessionID: sessionID,
		Identity: &identity.Identity{
			Emails:   []identity.Email{{Value: email}},
			Timezone: timezone,
		},
		Password: password,
	}
	ctx := context.Background()
	// Mocking expected behavior for valid flow and session
	s.mockFlowRepository.On("QueryFlowByFlowID", ctx, registrationFlow).
		Run(func(args mock.Arguments) {
			registrationFlow := args.Get(1).(*registration.Flow)
			registrationFlow.ExpiresAt = time.Time{} // Expired flow
		}).
		Return(nil).Once()
	var err error
	name := "registrationFlows/" + uuid.New().String()
	_, err = s.service.CompleteRegistrationFlow(ctx, registrationFlow, name, ipAddress, userAgent)
	// Assert: Ensure no error and valid session returned
	s.Require().Equal(ErrFlowExpired.Error(), err.Error())

	s.mockFlowRepository.AssertExpectations(s.T())
	s.mockSessionRepository.AssertExpectations(s.T())
	s.mockIdentityRepository.AssertExpectations(s.T())
}

func (s *serviceTestSuite) TestCompleteRegistrationFlow_SessionExpire() {
	ipAddress, _ := netip.ParseAddr("192.168.1.1")
	// Arrange: create a valid flow and session
	registrationFlow := &registration.Flow{
		SessionID: sessionID,
		Identity: &identity.Identity{
			Emails:   []identity.Email{{Value: email}},
			Timezone: timezone,
		},
		Password: password,
	}
	ctx := context.Background()
	// Mocking expected behavior for valid flow and session
	s.mockFlowRepository.On("QueryFlowByFlowID", ctx, registrationFlow).
		Run(func(args mock.Arguments) {
			registrationFlow := args.Get(1).(*registration.Flow)
			registrationFlow.ExpiresAt = time.Now().Add(900 * time.Hour)
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("QuerySessionWithDevices", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			preSession := args.Get(1).(*session.Session)
			preSession.Devices = []session.Device{
				{GeoLocation: "not the same anyway"},
			}
			preSession.ExpiresAt = time.Time{} // Expired session
		}).
		Return(nil).Once()

	var err error
	name := "registrationFlows/" + uuid.New().String()
	_, err = s.service.CompleteRegistrationFlow(ctx, registrationFlow, name, ipAddress, userAgent)
	// Assert: Ensure no error and valid session returned
	s.Require().Equal(ErrSessionExpired.Error(), err.Error())
	// Assert: Ensure mocks were called
	s.mockFlowRepository.AssertExpectations(s.T())
	s.mockSessionRepository.AssertExpectations(s.T())
	s.mockIdentityRepository.AssertExpectations(s.T())
}

func (s *serviceTestSuite) TestCompleteRegistrationFlow_DeviceExist() {
	ipAddress, _ := netip.ParseAddr("192.168.1.1")
	// Arrange: create a valid flow and session
	registrationFlow := &registration.Flow{
		SessionID: sessionID,
		Identity: &identity.Identity{
			Emails:   []identity.Email{{Value: email}},
			Timezone: timezone,
		},
		Password: password,
	}

	ctx := context.Background()
	// Mocking expected behavior for valid flow and session
	s.mockFlowRepository.On("QueryFlowByFlowID", ctx, registrationFlow).
		Run(func(args mock.Arguments) {
			registrationFlow := args.Get(1).(*registration.Flow)
			registrationFlow.ExpiresAt = time.Now().Add(900 * time.Hour)
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("QuerySessionWithDevices", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			preSession := args.Get(1).(*session.Session)
			preSession.ExpiresAt = time.Now().Add(900 * time.Hour)
			preSession.Devices = []session.Device{
				{
					IPAddress:   ipAddress,
					UserAgent:   userAgent,
					GeoLocation: "(unimplemented)",
				},
			}
		}).
		Return(nil).Once()
	s.mockIdentityRepository.On("EmailExists", ctx, registrationFlow.Identity.Emails[0].Value).Return(false, nil).Once() // Email doesn't exist
	s.mockIdentityRepository.On("CreateIdentity", ctx, mock.Anything).Return(nil).Once()
	s.mockSessionRepository.On("CreateSession", ctx, mock.Anything).Return(nil).Once()
	// Act: call CompleteRegistrationFlow
	name := "registrationFlows/" + uuid.New().String()
	_, err := s.service.CompleteRegistrationFlow(ctx, registrationFlow, name, ipAddress, userAgent)
	s.Require().NoError(err)
	// Assert: Ensure no error and valid session returned
	// Assert: Ensure mocks were called
	s.mockFlowRepository.AssertExpectations(s.T())
	s.mockSessionRepository.AssertExpectations(s.T())
	s.mockIdentityRepository.AssertExpectations(s.T())
}

func (s *serviceTestSuite) TestCompleteRegistrationFlow_EmailExists() {
	ipAddress, _ := netip.ParseAddr("192.168.1.1")
	// Arrange: create a valid flow and session
	registrationFlow := &registration.Flow{
		SessionID: sessionID,
		Identity: &identity.Identity{
			Emails:   []identity.Email{{Value: email}},
			Timezone: timezone,
		},
		Password: password,
	}
	ctx := context.Background()
	// Mocking expected behavior for valid flow and session
	s.mockFlowRepository.On("QueryFlowByFlowID", ctx, registrationFlow).
		Run(func(args mock.Arguments) {
			registrationFlow := args.Get(1).(*registration.Flow)
			registrationFlow.ExpiresAt = time.Now().Add(900 * time.Hour)
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("QuerySessionWithDevices", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			preSession := args.Get(1).(*session.Session)
			preSession.Devices = []session.Device{
				{GeoLocation: "not the same anyway"},
			}
			preSession.ExpiresAt = time.Now().Add(900 * time.Hour)
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("InsertDevice", ctx, mock.Anything).Return(nil).Once()
	s.mockIdentityRepository.On("EmailExists", ctx, registrationFlow.Identity.Emails[0].Value).Return(true, nil).Once() // Email exist

	// Act: call CompleteRegistrationFlow
	var err error
	name := "registrationFlows/" + uuid.New().String()
	_, err = s.service.CompleteRegistrationFlow(ctx, registrationFlow, name, ipAddress, userAgent)
	// Assert: Ensure no error and valid session returned
	s.Require().Equal(ErrEmailExists.Error(), err.Error())
	// Assert: Ensure mocks were called
	s.mockFlowRepository.AssertExpectations(s.T())
	s.mockSessionRepository.AssertExpectations(s.T())
	s.mockIdentityRepository.AssertExpectations(s.T())
}

func (s *serviceTestSuite) TestCompleteRegistrationFlow_WeakPassword() {
	ipAddress, _ := netip.ParseAddr("192.168.1.1")
	// Arrange: create a valid flow and session
	registrationFlow := &registration.Flow{
		SessionID: sessionID,
		Identity: &identity.Identity{
			Emails:   []identity.Email{{Value: email}},
			Timezone: timezone,
		},
		Password: "NoNumberPassword",
	}
	ctx := context.Background()
	// Mocking expected behavior for valid flow and session
	s.mockFlowRepository.On("QueryFlowByFlowID", ctx, registrationFlow).
		Run(func(args mock.Arguments) {
			registrationFlow := args.Get(1).(*registration.Flow)
			registrationFlow.ExpiresAt = time.Now().Add(900 * time.Hour)
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("QuerySessionWithDevices", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			preSession := args.Get(1).(*session.Session)
			preSession.Devices = []session.Device{
				{GeoLocation: "not the same anyway"},
			}
			preSession.ExpiresAt = time.Now().Add(900 * time.Hour)
		}).
		Return(nil).Once()
	s.mockSessionRepository.On("InsertDevice", ctx, mock.Anything).Return(nil).Once()
	s.mockIdentityRepository.On("EmailExists", ctx, registrationFlow.Identity.Emails[0].Value).Return(false, nil).Once() // Email doesn't exist
	// Act: call CompleteRegistrationFlow
	var err error
	name := "registrationFlows/" + uuid.New().String()
	_, err = s.service.CompleteRegistrationFlow(ctx, registrationFlow, name, ipAddress, userAgent)
	// Assert: Ensure no error and valid session returned
	s.Require().Equal(ErrInsecurePassword.Error(), err.Error())

	// Assert: Ensure mocks were called
	s.mockFlowRepository.AssertExpectations(s.T())
	s.mockSessionRepository.AssertExpectations(s.T())
	s.mockIdentityRepository.AssertExpectations(s.T())
}

func (s *serviceTestSuite) TestCompleteRegistrationFlow_NoNameInFlow() {
	ipAddress, _ := netip.ParseAddr("192.168.1.1")
	// Arrange: create a valid flow and session
	registrationFlow := &registration.Flow{
		SessionID: sessionID,
		Identity: &identity.Identity{
			Emails:   []identity.Email{{Value: email}},
			Timezone: timezone,
		},
		Password: "NoNumberPassword",
	}
	ctx := context.Background()
	// Mocking expected behavior for valid flow and session
	var err error
	_, err = s.service.CompleteRegistrationFlow(ctx, registrationFlow, "", ipAddress, userAgent)
	// Assert: Ensure no error and valid session returned
	s.Require().Equal(ErrUnauthenticated.Error(), err.Error())

	// Assert: Ensure mocks were called
	s.mockFlowRepository.AssertExpectations(s.T())
	s.mockSessionRepository.AssertExpectations(s.T())
	s.mockIdentityRepository.AssertExpectations(s.T())
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}
