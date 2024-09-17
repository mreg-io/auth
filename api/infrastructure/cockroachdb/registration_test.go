package cockroachdb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"gitlab.mreg.io/my-registry/auth/domain/registration"
)

type RegistrationRepositorySuite struct {
	suite.Suite
	pool       *pgxpool.Pool
	repository registration.Repository
}

func (s *RegistrationRepositorySuite) SetupSuite() {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	s.Require().NoError(err)
	s.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	s.Require().NoError(err)
	s.repository = NewRegistrationRepository(s.pool)
	s.T().Setenv("REGISTRATION_EXPIRY_INTERVAL", "1h")
	s.T().Setenv("SESSION_EXPIRY_INTERVAL", "1h")
}

func (s *RegistrationRepositorySuite) TestCreateFlow_WithoutErr() {
	ctx := context.Background()
	SessionInterval, err := time.ParseDuration(os.Getenv("SESSION_EXPIRY_INTERVAL"))
	s.Require().NoError(err)
	registrationInterval, err := time.ParseDuration(os.Getenv("REGISTRATION_EXPIRY_INTERVAL"))
	s.Require().NoError(err)
	var savedSessionID string
	var savedExpiry time.Time
	err = s.pool.
		QueryRow(ctx, `INSERT INTO sessions (active, authenticator_assurance_level, expires_at)
        VALUES (true, 1, current_timestamp + $1)RETURNING id, expires_at`, SessionInterval).
		Scan(&savedSessionID, &savedExpiry)
	s.Require().NoError(err)
	testCase := &registration.Flow{
		Interval:  registrationInterval,
		SessionID: savedSessionID,
	}

	err = s.repository.CreateFlow(ctx, testCase)
	s.Require().NoError(err)

	var savedFlow registration.Flow
	err = s.pool.
		QueryRow(ctx, `SELECT id, issued_at, expires_at, session_id from registration_flows WHERE id = $1`, testCase.FlowID).
		Scan(&savedFlow.FlowID, &savedFlow.IssuedAt, &savedFlow.ExpiresAt, &savedFlow.SessionID)
	s.Require().NoError(err)
	s.Equal(testCase.FlowID, savedFlow.FlowID)
	s.Equal(savedFlow.IssuedAt, testCase.IssuedAt)
	s.Equal(savedFlow.ExpiresAt, testCase.ExpiresAt)
	s.Greater(savedFlow.ExpiresAt, savedFlow.IssuedAt)
}

func (s *RegistrationRepositorySuite) TestCreateFlow_WithUnknownSessionID() {
	ctx := context.Background()
	SessionInterval, err := time.ParseDuration(os.Getenv("SESSION_EXPIRY_INTERVAL"))
	s.Require().NoError(err)
	registrationInterval, err := time.ParseDuration(os.Getenv("REGISTRATION_EXPIRY_INTERVAL"))
	s.Require().NoError(err)
	var savedSessionID string
	var savedExpiry time.Time
	err = s.pool.
		QueryRow(ctx, `INSERT INTO sessions (active, authenticator_assurance_level, expires_at)
        VALUES (true, 1, current_timestamp + $1)RETURNING id, expires_at`, SessionInterval).
		Scan(&savedSessionID, &savedExpiry)
	s.Require().NoError(err)
	testCase := &registration.Flow{
		Interval:  registrationInterval,
		SessionID: "53346325-ed55-4ab1-8ca5-ed5bcad352a1",
	}

	err = s.repository.CreateFlow(ctx, testCase)
	s.Require().Error(err)

	var savedFlow registration.Flow
	err = s.pool.
		QueryRow(ctx, `SELECT id, issued_at, expires_at, session_id from registration_flows WHERE id = $1`, testCase.FlowID).
		Scan(&savedFlow.FlowID, &savedFlow.IssuedAt, &savedFlow.ExpiresAt, &savedFlow.SessionID)
	s.Require().Error(err)
}

func (s *RegistrationRepositorySuite) TearDownSuite() {
	s.pool.Close()
}

func TestRegistrationRepositorySuite(t *testing.T) {
	suite.Run(t, new(RegistrationRepositorySuite))
}
