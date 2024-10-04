package cockroachdb

import (
	"context"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"gitlab.mreg.io/my-registry/auth/domain/session"
)

type SessionRepositorySuite struct {
	suite.Suite
	pool       *pgxpool.Pool
	repository session.Repository
}

var (
	sessionSessionID1  = uuid.New() // for query
	sessionSessionID2  = uuid.New() // for query
	sessionSessionID3  = uuid.New() // for query
	sessionSessionID4  = uuid.New() // for update device
	sessionIdentityID1 = uuid.New() // for query
	sessionIdentityID2 = uuid.New() // for query
	sessionDeviceID1   = uuid.New() // for query
	sessionDeviceID2   = uuid.New() // for query
	sessionDeviceID3   = uuid.New() // for query
	sessionDeviceID4   = uuid.New() // for update device
)

func (s *SessionRepositorySuite) SetupSuite() {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	s.Require().NoError(err)
	s.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	s.Require().NoError(err)
	s.repository = NewSessionRepository(s.pool)

	ctx := context.Background()
	_, err = s.pool.Exec(ctx, `
        INSERT INTO identities (id, timezone) 
        VALUES 
        ($1, 'Thailand' ),
        ($2, 'Thailand' )
    `, sessionIdentityID1, sessionIdentityID2)
	s.Require().NoError(err)

	_, err = s.pool.Exec(ctx, `
        INSERT INTO sessions (id, active, authenticator_assurance_level,
                              issued_at, expires_at, identity_id) 
        VALUES 
        ($1, true, 1, current_timestamp,current_timestamp + interval '2 hours', $5),
        ($2, false, NULL, current_timestamp, current_timestamp + interval '2 hours', $6),
        ($3, true, 2, current_timestamp, current_timestamp + interval '2 hours', NULL),
        ($4, true, 2, current_timestamp, current_timestamp + interval '2 hours', NULL)
    `, sessionSessionID1, sessionSessionID2, sessionSessionID3, sessionSessionID4,
		sessionIdentityID1, sessionIdentityID2)
	s.Require().NoError(err)

	// session2 has two devices
	_, err = s.pool.Exec(ctx, `
        INSERT INTO devices (id, ip_address,geo_location, user_agent, session_id) 
        VALUES 
        ($1, '192.168.0.1','Thailand', 'Mozilla', $5),
        ($2, '192.168.1.1','Thailand', 'Gorilla', $6),
        ($3, '192.168.1.2','Taiwan', 'UA', $7),
        ($4, '192.168.87.87','InsertDevice', 'BunBun', $8)
    `, sessionDeviceID1, sessionDeviceID2, sessionDeviceID3, sessionDeviceID4,
		sessionSessionID1, sessionSessionID2, sessionSessionID2, sessionSessionID4)
	s.Require().NoError(err)
}

func (s *SessionRepositorySuite) TestCreateSession_WithoutErr() {
	ctx := context.Background()
	interval, err := time.ParseDuration("2h")
	s.Require().NoError(err)
	addr, err := netip.ParseAddr("118.232.60.138")
	s.Require().NoError(err)

	testCase := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 1,
		Devices: []session.Device{
			{
				IPAddress:   addr,
				GeoLocation: "TW",
				UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.15",
			},
		},
		ExpiryInterval: interval,
	}

	err = s.repository.CreateSession(ctx, testCase)
	s.Require().NoError(err)

	var savedSession session.Session
	err = s.pool.
		QueryRow(ctx, `SELECT id, active, coalesce(authenticator_assurance_level, 0), issued_at, expires_at from sessions WHERE id = $1`, testCase.ID).
		Scan(&savedSession.ID, &savedSession.Active, &savedSession.AuthenticatorAssuranceLevel, &savedSession.IssuedAt, &savedSession.ExpiresAt)
	s.Require().NoError(err)
	s.Equal(testCase.ID, savedSession.ID)
	s.True(savedSession.Active)
	s.Equal(uint8(1), savedSession.AuthenticatorAssuranceLevel)
	s.Equal(savedSession.IssuedAt, testCase.IssuedAt)
	s.Equal(savedSession.ExpiresAt, testCase.IssuedAt.Add(interval))
	s.Equal(savedSession.AuthenticatedAt, testCase.AuthenticatedAt)
}

func (s *SessionRepositorySuite) TestCreateSession_0AALWithoutErr() {
	ctx := context.Background()
	interval, err := time.ParseDuration("2h")
	s.Require().NoError(err)
	addr, err := netip.ParseAddr("118.232.60.138")
	s.Require().NoError(err)

	testCase := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 0,
		Devices: []session.Device{
			{
				IPAddress:   addr,
				GeoLocation: "TW",
				UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.15",
			},
		},
		ExpiryInterval: interval,
	}

	err = s.repository.CreateSession(ctx, testCase)
	s.Require().NoError(err)

	var savedSession session.Session
	err = s.pool.
		QueryRow(ctx, `SELECT id, active, COALESCE(authenticator_assurance_level, 0), issued_at, expires_at from sessions WHERE id = $1`, testCase.ID).
		Scan(&savedSession.ID, &savedSession.Active, &savedSession.AuthenticatorAssuranceLevel, &savedSession.IssuedAt, &savedSession.ExpiresAt)
	s.Require().NoError(err)
	s.Equal(testCase.ID, savedSession.ID)
	s.True(true, savedSession.Active)
	s.Equal(uint8(0), testCase.AuthenticatorAssuranceLevel)
	s.Equal(savedSession.IssuedAt, testCase.IssuedAt)
	s.Equal(savedSession.ExpiresAt, testCase.IssuedAt.Add(interval))
	s.Equal(savedSession.AuthenticatedAt, testCase.AuthenticatedAt)
}

func (s *SessionRepositorySuite) TestCreateSession_EmptyDevice() {
	ctx := context.Background()
	interval, err := time.ParseDuration("2h")
	s.Require().NoError(err)

	testCase := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 1,
		ExpiryInterval:              interval,
	}

	err = s.repository.CreateSession(ctx, testCase)
	s.Require().Error(err)
}

func (s *SessionRepositorySuite) TestCreateSession_MultipleDevices() {
	ctx := context.Background()
	interval, err := time.ParseDuration("2h")
	s.Require().NoError(err)
	addr1, err := netip.ParseAddr("221.232.60.138")
	s.Require().NoError(err)
	addr2, err := netip.ParseAddr("118.232.60.138")
	s.Require().NoError(err)

	testCase := &session.Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 1,
		Devices: []session.Device{
			{
				IPAddress:   addr1,
				GeoLocation: "TW",
				UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.15",
			},
			{
				IPAddress:   addr2,
				GeoLocation: "TW",
				UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36 Edg/97.0.1072.71",
			},
		},
		ExpiryInterval: interval,
	}

	err = s.repository.CreateSession(ctx, testCase)
	s.Require().Error(err)
}

func (s *SessionRepositorySuite) TestQuerySessionByID_NoErr() {
	ctx := context.Background()

	// below is the test
	sessionData := &session.Session{
		ID: sessionSessionID1.String(),
	}
	err := s.repository.QuerySessionByID(ctx, sessionData)
	s.Require().NoError(err)
	s.Require().True(sessionData.Active)
	s.Require().Equal(uint8(1), sessionData.AuthenticatorAssuranceLevel)
	s.Require().NotEmpty(sessionData.IssuedAt)
	s.Require().NotEmpty(sessionData.ExpiresAt)
	s.Require().Equal(sessionIdentityID1.String(), sessionData.Identity.ID)
}

func (s *SessionRepositorySuite) TestQuerySessionByID_NotExistSessionID() {
	ctx := context.Background()

	// below is the test
	sessionData := &session.Session{
		ID: uuid.New().String(),
	}
	err := s.repository.QuerySessionByID(ctx, sessionData)
	s.Require().Error(err)
}

func (s *SessionRepositorySuite) TestQuerySessionByID_AuthenticatorAssuranceLevel_Is_NULL_NoErr() {
	ctx := context.Background()

	// below is the test
	sessionData := &session.Session{
		ID: sessionSessionID2.String(),
	}
	err := s.repository.QuerySessionByID(ctx, sessionData)
	s.Require().NoError(err)
	s.Require().False(sessionData.Active)
	s.Require().Equal(uint8(0), sessionData.AuthenticatorAssuranceLevel)
	s.Require().NotEmpty(sessionData.IssuedAt)
	s.Require().NotEmpty(sessionData.ExpiresAt)
	s.Require().Equal(sessionIdentityID2.String(), sessionData.Identity.ID)
}

func (s *SessionRepositorySuite) TestQuerySessionByID_NoIdentity_NoErr() {
	ctx := context.Background()

	// below is the test
	sessionData := &session.Session{
		ID: sessionSessionID3.String(),
	}
	err := s.repository.QuerySessionByID(ctx, sessionData)
	s.Require().NoError(err)
	s.Require().True(sessionData.Active)
	s.Require().Equal(uint8(2), sessionData.AuthenticatorAssuranceLevel)
	s.Require().NotEmpty(sessionData.IssuedAt)
	s.Require().NotEmpty(sessionData.ExpiresAt)
	s.Require().Zero(sessionData.Identity.ID)
}

func (s *SessionRepositorySuite) TestQuerySessionWithDevice_OneDevice_NoErr() {
	ctx := context.Background()

	// below is the test
	sessionData := &session.Session{
		ID: sessionSessionID1.String(),
	}
	err := s.repository.QuerySessionWithDevices(ctx, sessionData)
	s.Require().NoError(err)
	s.Require().True(sessionData.Active)
	s.Require().Equal(uint8(1), sessionData.AuthenticatorAssuranceLevel)
	s.Require().NotEmpty(sessionData.IssuedAt)
	s.Require().NotEmpty(sessionData.ExpiresAt)
	s.Require().Equal(sessionIdentityID1.String(), sessionData.Identity.ID)

	// device related
	s.Require().Len(sessionData.Devices, 1)
	s.Require().Equal(sessionSessionID1.String(), sessionData.Devices[0].SessionID)
	s.Require().Equal("192.168.0.1", sessionData.Devices[0].IPAddress.String())
	s.Require().Equal("Mozilla", sessionData.Devices[0].UserAgent)
	s.Require().Equal("Thailand", sessionData.Devices[0].GeoLocation)
}

func (s *SessionRepositorySuite) TestQuerySessionWithDevice_TwoDevice_NoErr() {
	ctx := context.Background()

	// below is the test
	sessionData := &session.Session{
		ID: sessionSessionID2.String(),
	}
	err := s.repository.QuerySessionWithDevices(ctx, sessionData)
	s.Require().NoError(err)
	s.Require().False(sessionData.Active)
	s.Require().Equal(uint8(0), sessionData.AuthenticatorAssuranceLevel)
	s.Require().NotEmpty(sessionData.IssuedAt)
	s.Require().NotEmpty(sessionData.ExpiresAt)
	s.Require().Equal(sessionIdentityID2.String(), sessionData.Identity.ID)

	s.Require().Len(sessionData.Devices, 2)
	s.Require().Equal(sessionSessionID2.String(), sessionData.Devices[0].SessionID)

	// make sure sequence is correct
	if sessionData.Devices[0].IPAddress.String() == "192.168.1.2" {
		sessionData.Devices[0], sessionData.Devices[1] = sessionData.Devices[1], sessionData.Devices[0]
	}
	s.Require().Equal("192.168.1.1", sessionData.Devices[0].IPAddress.String())
	s.Require().Equal("Thailand", sessionData.Devices[0].GeoLocation)
	s.Require().Equal("Gorilla", sessionData.Devices[0].UserAgent)

	s.Require().Equal(sessionSessionID2.String(), sessionData.Devices[1].SessionID)
	s.Require().Equal("192.168.1.2", sessionData.Devices[1].IPAddress.String())
	s.Require().Equal("Taiwan", sessionData.Devices[1].GeoLocation)
	s.Require().Equal("UA", sessionData.Devices[1].UserAgent)
}

func (s *SessionRepositorySuite) TestQuerySessionWithDevice_NotExistSessionID_Err() {
	ctx := context.Background()

	// below is the test
	sessionData := &session.Session{
		ID: uuid.New().String(),
	}
	err := s.repository.QuerySessionWithDevices(ctx, sessionData)
	s.Require().Error(err)
}

func (s *SessionRepositorySuite) TestUpdateDevice_NoErr() {
	ctx := context.Background()

	// below is the test
	device := &session.Device{
		SessionID:   sessionSessionID4.String(),
		IPAddress:   netip.MustParseAddr("192.168.69.2"),
		UserAgent:   "Chrome",
		GeoLocation: "USA",
	}
	err := s.repository.InsertDevice(ctx, device)
	s.Require().NoError(err)

	sessionData := &session.Session{
		ID: sessionSessionID4.String(),
	}
	err = s.repository.QuerySessionWithDevices(ctx, sessionData)

	// make sure sessionData.Devices[1] is the new one
	if sessionData.Devices[1].IPAddress.String() == "192.168.87.87" {
		sessionData.Devices[0], sessionData.Devices[1] = sessionData.Devices[1], sessionData.Devices[0]
	}

	s.Require().NoError(err)
	s.Require().Len(sessionData.Devices, 2)
	s.Require().NotEmpty(sessionData.Devices[1].ID)
	s.Require().Equal(sessionSessionID4.String(), sessionData.Devices[1].SessionID)
	s.Require().Equal("192.168.69.2", sessionData.Devices[1].IPAddress.String())
	s.Require().Equal("Chrome", sessionData.Devices[1].UserAgent)
	s.Require().Equal("USA", sessionData.Devices[1].GeoLocation)
}

func (s *SessionRepositorySuite) TestUpdateDevice_NotExistSession_Err() {
	ctx := context.Background()

	// below is the test
	device := &session.Device{
		SessionID:   uuid.New().String(),
		IPAddress:   netip.MustParseAddr("192.168.69.2"),
		UserAgent:   "Chrome",
		GeoLocation: "USA",
	}
	err := s.repository.InsertDevice(ctx, device)
	s.Require().Error(err)
}

func (s *SessionRepositorySuite) TestUpdateDevice_NotProvideSessionID_Err() {
	ctx := context.Background()

	// below is the test
	device := &session.Device{
		IPAddress:   netip.MustParseAddr("192.168.69.2"),
		UserAgent:   "Chrome",
		GeoLocation: "USA",
	}
	err := s.repository.InsertDevice(ctx, device)
	s.Require().Error(err)
}

func (s *SessionRepositorySuite) TearDownSuite() {
	s.pool.Close()
}

func TestSessionRepositorySuite(t *testing.T) {
	suite.Run(t, new(SessionRepositorySuite))
}
