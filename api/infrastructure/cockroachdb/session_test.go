package cockroachdb

import (
	"context"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"gitlab.mreg.io/my-registry/auth/domain/session"
)

type SessionRepositorySuite struct {
	suite.Suite
	pool       *pgxpool.Pool
	repository session.Repository
}

func (s *SessionRepositorySuite) SetupSuite() {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	s.Require().NoError(err)
	s.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	s.Require().NoError(err)
	s.repository = NewSessionRepository(s.pool)
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

func (s *SessionRepositorySuite) TearDownSuite() {
	s.pool.Close()
}

func TestSessionRepositorySuite(t *testing.T) {
	suite.Run(t, new(SessionRepositorySuite))
}
