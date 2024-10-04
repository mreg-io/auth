package cockroachdb

import (
	"context"
	"os"
	"testing"

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
	identityEmail1 = generateRandomEmail()
	password1      = "password123"
)

func (i *IdentityRepositorySuite) SetupSuite() {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	i.Require().NoError(err)
	i.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	i.Require().NoError(err)
	i.repository = NewIdentityRepository(i.pool)
}

func (i *IdentityRepositorySuite) TestCreateIdentity_WithoutErr() {
	ctx := context.Background()
	var err error
	newEmail := identity.Email{
		Value:    identityEmail1,
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

func (i *IdentityRepositorySuite) TearDownSuite() {
	i.pool.Close()
}

func TestIdentityRepositorySuite(t *testing.T) {
	suite.Run(t, new(IdentityRepositorySuite))
}
