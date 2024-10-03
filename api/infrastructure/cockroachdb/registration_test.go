package cockroachdb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.mreg.io/my-registry/auth/domain/identity"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"gitlab.mreg.io/my-registry/auth/domain/registration"
)

type RegistrationRepositorySuite struct {
	suite.Suite
	pool       *pgxpool.Pool
	repository registration.Repository
}

var (
	registrationSessionID1 = uuid.New()
	registrationSessionID2 = uuid.New()
	registrationSessionID3 = uuid.New()
	registrationFlowID1    = uuid.New()
	registrationFlowID2    = uuid.New()
	registrationFlowID3    = uuid.New()
)

func (s *RegistrationRepositorySuite) SetupSuite() {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	s.Require().NoError(err)
	s.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	s.Require().NoError(err)
	s.repository = NewRegistrationRepository(s.pool)
	s.T().Setenv("REGISTRATION_EXPIRY_INTERVAL", "1h")
	s.T().Setenv("SESSION_EXPIRY_INTERVAL", "1h")

	ctx := context.Background()
	_, err = s.pool.Exec(ctx, `
        INSERT INTO sessions (id, active, authenticator_assurance_level, issued_at, expires_at) 
        VALUES 
        ($1, true, 1, current_timestamp, current_timestamp + interval '2 hours'),
        ($2, false, 1, current_timestamp, current_timestamp + interval '2 hours'),
        ($3, true, 2, current_timestamp, current_timestamp + interval '2 hours')
    `, registrationSessionID1, registrationSessionID2, registrationSessionID3)
	s.Require().NoError(err)

	_, err = s.pool.Exec(ctx, `
        INSERT INTO registration_flows (id, issued_at, expires_at, session_id) 
        VALUES 
        ($1, current_timestamp, current_timestamp + interval '2 hours', $4),
        ($2, current_timestamp, current_timestamp + interval '2 hours', $5),
        ($3, current_timestamp, current_timestamp + interval '2 hours', $6)
    `, registrationFlowID1, registrationFlowID2, registrationFlowID3,
		registrationSessionID1, registrationSessionID2, registrationSessionID3)
	s.Require().NoError(err)
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

func (s *RegistrationRepositorySuite) TestQueryFlow1() {
	ctx := context.Background()

	// below is the test
	flow := &registration.Flow{
		FlowID: registrationFlowID1.String(),
	}
	err := s.repository.QueryFlowByFlowID(ctx, flow)
	s.Require().NoError(err)
	s.Require().Equal(registrationSessionID1.String(), flow.SessionID)
	s.Require().NotEmpty(flow.IssuedAt)
	s.Require().NotEmpty(flow.ExpiresAt)
	s.Greater(flow.ExpiresAt, flow.IssuedAt)
}

func (s *RegistrationRepositorySuite) TestQueryFlow2_SetUpAllField_ShouldOverWrite() {
	issuedAt, err := time.Parse(time.UnixDate, "Wed Feb 25 11:06:39 PST 5069")
	s.Require().NoError(err)
	expiresAt, err := time.Parse(time.UnixDate, "Wed Feb 25 11:06:39 PST 2069")
	s.Require().NoError(err)
	ctx := context.Background()

	flow := &registration.Flow{
		FlowID:    registrationFlowID2.String(),
		SessionID: "ed074d1e-fe04-4683-9239-91cf59f126c9",
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Password:  "secret",
		Interval:  time.Hour,
		Identity:  &identity.Identity{},
	}
	err = s.repository.QueryFlowByFlowID(ctx, flow)
	s.Require().NoError(err)
	s.Require().NotEqual("ed074d1e-fe04-4683-9239-91cf59f126c9", flow.SessionID)
	s.Require().NotEqual(issuedAt, flow.IssuedAt)
	s.Require().NotEqual(expiresAt, flow.ExpiresAt)
	s.Greater(flow.ExpiresAt, flow.IssuedAt)
}

func (s *RegistrationRepositorySuite) TestQueryFlow3_SearchNotExistFlowID() {
	ctx := context.Background()

	// below is the test
	flow := &registration.Flow{
		FlowID: uuid.New().String(),
	}
	err := s.repository.QueryFlowByFlowID(ctx, flow)
	s.Require().Error(err)
}

func (s *RegistrationRepositorySuite) TearDownSuite() {
	s.pool.Close()
}

func TestRegistrationRepositorySuite(t *testing.T) {
	suite.Run(t, new(RegistrationRepositorySuite))
}
