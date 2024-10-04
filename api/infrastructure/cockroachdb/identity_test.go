package cockroachdb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"gitlab.mreg.io/my-registry/auth/domain/identity"
)

type IdentityRepositorySuite struct {
	suite.Suite
	pool       *pgxpool.Pool
	repository identity.Repository
}

func generateRandomEmail() string {
	localPart := "User" + uuid.New().String() // Generates a random number from 0 to 999
	domain := "example.com"                   // Domain name
	return localPart + "@" + domain
}

var (
	identityIdentityID1            = uuid.New() // for query
	identityIdentityID2            = uuid.New() // for query
	identityEmail1                 = generateRandomEmail()
	identityEmail2                 = generateRandomEmail()
	verifiedAt1                    = time.Date(6069, time.October, 4, 5, 35, 44, 155200000, time.UTC)
	identityEmailForCreateIdentity = generateRandomEmail()
	password1                      = "password123"
)

func (i *IdentityRepositorySuite) SetupSuite() {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	i.Require().NoError(err)
	i.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	i.Require().NoError(err)
	i.repository = NewIdentityRepository(i.pool)

	ctx := context.Background()
	_, err = i.pool.Exec(ctx, `
        INSERT INTO identities (id, timezone) 
        VALUES 
        ($1, 'Thailand'),
        ($2, 'Taiwan')
    `, identityIdentityID1, identityIdentityID2)
	i.Require().NoError(err)

	_, err = i.pool.Exec(ctx, `
        INSERT INTO emails (address, verified, verified_at, identity_id) 
        VALUES 
        ($1, DEFAULT, NULL, $2),
        ($3, True, $4, $5)
    `, identityEmail1, identityIdentityID1,
		identityEmail2, verifiedAt1, identityIdentityID2)
	i.Require().NoError(err)
}

func (i *IdentityRepositorySuite) TestCreateIdentity_WithoutErr() {
	ctx := context.Background()
	var err error
	newEmail := identity.Email{
		Value:    identityEmailForCreateIdentity,
		Verified: false,
	}
	newIdentity := &identity.Identity{
		Emails:       []identity.Email{newEmail},
		PasswordHash: password1,
		Timezone:     "Taipei/Taiwan",
		State:        identity.StateActive,
	}
	originalIdentity := *newIdentity
	err = i.repository.CreateIdentity(ctx, newIdentity)
	i.Require().NoError(err)

	// check the data integrity of the newIdentity
	i.Require().Equal(originalIdentity.Emails[0].Value, newIdentity.Emails[0].Value)
	i.Require().Equal(originalIdentity.Emails[0].Verified, newIdentity.Emails[0].Verified)
	i.Require().Equal(originalIdentity.PasswordHash, newIdentity.PasswordHash)
	i.Require().Equal(originalIdentity.Timezone, newIdentity.Timezone)
	i.Require().Equal(originalIdentity.State, newIdentity.State)

	dbIdentity := identity.Identity{}
	err = i.pool.
		QueryRow(ctx, `SELECT id, 
       	CASE state 
        	WHEN 'active' THEN 1 
        	WHEN 'suspended' THEN 2 
    	END AS state, 
    	timezone, create_time, update_time, state_update_time 
		FROM identities WHERE id=$1`, newIdentity.ID).
		Scan(&dbIdentity.ID, &dbIdentity.State, &dbIdentity.Timezone,
			&dbIdentity.CreateTime, &dbIdentity.UpdateTime, &dbIdentity.StateUpdateTime)
	i.Require().NoError(err)

	dbEmail := identity.Email{}
	err = i.pool.
		QueryRow(ctx, `SELECT address, verified, create_time, update_time FROM emails WHERE identity_id=$1`, newIdentity.ID).
		Scan(&dbEmail.Value, &dbEmail.Verified, &dbEmail.CreateTime, &dbEmail.UpdateTime)
	i.Require().NoError(err)

	var dbPassword string
	var dbPasswordIdentityID string
	err = i.pool.
		QueryRow(ctx, `SELECT * FROM passwords WHERE identity_id=$1`, newIdentity.ID).
		Scan(&dbPasswordIdentityID, &dbPassword)
	i.Require().NoError(err)

	// test for identity table
	i.Require().NotEmpty(dbIdentity.ID, "Identity ID should not be empty")
	i.Require().Equal(identity.StateActive, dbIdentity.State, "Identity state should be active")
	i.Require().Equal(newIdentity.Timezone, dbIdentity.Timezone, "Identity timezone should match")
	i.Require().NotZero(dbIdentity.CreateTime, "CreateTime should not be zero")
	i.Require().NotZero(dbIdentity.UpdateTime, "UpdateTime should not be zero")
	i.Require().NotZero(dbIdentity.StateUpdateTime, "StateUpdateTime should not be zero")

	// test for email table
	i.Require().Equal(newEmail.Value, dbEmail.Value, "Email address should match")
	i.Require().Equal(newEmail.Verified, dbEmail.Verified, "Email verified status should match")
	i.Require().NotZero(dbEmail.CreateTime, "Email CreateTime should not be zero")
	i.Require().NotZero(dbEmail.UpdateTime, "Email UpdateTime should not be zero")

	// test for password table
	i.Require().Equal(newIdentity.ID, dbPasswordIdentityID, "Password identity ID should match")
	i.Require().Equal(password1, dbPassword, "Password hash should match")
}

func (i *IdentityRepositorySuite) TestCreateIdentity_NoTimeZone_NoErr() {
	ctx := context.Background()
	var err error
	newEmail := identity.Email{
		Value:    generateRandomEmail(),
		Verified: false,
	}
	newIdentity := &identity.Identity{
		Emails:       []identity.Email{newEmail},
		PasswordHash: password1,
		State:        identity.StateActive,
	}
	err = i.repository.CreateIdentity(ctx, newIdentity)
	i.Require().NoError(err)
}

func (i *IdentityRepositorySuite) TestCreateIdentity_NoEmail_Err() {
	ctx := context.Background()
	var err error
	newIdentity := &identity.Identity{
		PasswordHash: password1,
		Timezone:     "Taipei/Taiwan",
		State:        identity.StateActive,
	}
	err = i.repository.CreateIdentity(ctx, newIdentity)
	i.Require().Error(err)
}

func (i *IdentityRepositorySuite) TestEmailExist() {
	ctx := context.Background()
	exist, err := i.repository.EmailExists(ctx, identityEmail1)
	i.Require().NoError(err)
	i.Require().True(exist)

	exist, err = i.repository.EmailExists(ctx, "non-existent-email@example.com")
	i.Require().NoError(err)
	i.Require().False(exist)
}

func (i *IdentityRepositorySuite) TestQueryEmail_NotVerifiedEmail_NoErr() {
	ctx := context.Background()
	queryEmail := &identity.Email{
		Value: identityEmail1,
	}
	err := i.repository.QueryEmail(ctx, queryEmail)
	i.Require().NoError(err)

	// check the Value not changed
	i.Require().Equal(identityEmail1, queryEmail.Value)

	i.Require().False(queryEmail.Verified)
	i.Require().NotEmpty(queryEmail.CreateTime)
	i.Require().NotEmpty(queryEmail.VerifiedAt) // doesn't really matter since not verified
	i.Require().NotEmpty(queryEmail.UpdateTime)
}

func (i *IdentityRepositorySuite) TestQueryEmail_VerifiedEmail_NoErr() {
	ctx := context.Background()
	queryEmail := &identity.Email{
		Value: identityEmail2,
	}
	err := i.repository.QueryEmail(ctx, queryEmail)
	i.Require().NoError(err)

	// check the Value not changed
	i.Require().Equal(identityEmail2, queryEmail.Value)

	i.Require().True(queryEmail.Verified)
	i.Require().NotEmpty(queryEmail.CreateTime)
	i.Require().Equal(verifiedAt1, queryEmail.VerifiedAt.UTC())
	i.Require().NotEmpty(queryEmail.UpdateTime)
}

func (i *IdentityRepositorySuite) TestQueryEmail_NotExistEmail_Err() {
	ctx := context.Background()
	queryEmail := &identity.Email{
		Value: generateRandomEmail(),
	}
	err := i.repository.QueryEmail(ctx, queryEmail)
	i.Require().Error(err)
}

func (i *IdentityRepositorySuite) TearDownSuite() {
	i.pool.Close()
}

func TestIdentityRepositorySuite(t *testing.T) {
	suite.Run(t, new(IdentityRepositorySuite))
}
