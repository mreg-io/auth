package cockroachdb

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.mreg.io/my-registry/auth/domain/identity"
)

type IdentityRepository struct {
	db *pgxpool.Pool
}

func NewIdentityRepository(db *pgxpool.Pool) identity.Repository {
	return &IdentityRepository{db: db}
}

//go:embed sql/createIdentity.sql
var createIdentitySQL string

//go:embed sql/queryEmail.sql
var queryEmailSQL string

func createIdentityField(identity *identity.Identity) []interface{} {
	return []interface{}{
		&identity.ID,
		&identity.CreateTime,
		&identity.UpdateTime,
		&identity.StateUpdateTime,
		&identity.Emails[0].CreateTime,
		&identity.Emails[0].UpdateTime,
	}
}

func (i *IdentityRepository) CreateIdentity(ctx context.Context, identityData *identity.Identity) error {
	if len(identityData.Emails) == 0 {
		return errors.New("identity must have at least one email")
	}
	return i.db.
		QueryRow(
			ctx,
			createIdentitySQL,
			identityData.Timezone, identityData.Emails[0].Value, identityData.PasswordHash,
		).
		Scan(createIdentityField(identityData)...)
}

func (i *IdentityRepository) EmailExists(ctx context.Context, emailAddress string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM emails WHERE address = $1)`

	var exists bool
	err := i.db.QueryRow(ctx, query, emailAddress).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func QueryEmailField(email *identity.Email) []interface{} {
	return []interface{}{
		&email.Verified,
		&email.CreateTime,
		&email.VerifiedAt,
		&email.UpdateTime,
	}
}

func (i *IdentityRepository) QueryEmail(ctx context.Context, email *identity.Email) error {
	return i.db.
		QueryRow(
			ctx,
			queryEmailSQL,
			email.Value,
		).
		Scan(QueryEmailField(email)...)
}
