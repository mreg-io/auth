package cockroachdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.mreg.io/my-registry/auth/domain/identity"
)

type IdentityRepository struct {
	db *pgxpool.Pool
}

func NewIdentityRepository(db *pgxpool.Pool) identity.Repository {
	return &IdentityRepository{db: db}
}

func (i *IdentityRepository) EmailExists(context.Context, string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (i *IdentityRepository) CreateIdentity(context.Context, *identity.Identity) error {
	// TODO implement me
	panic("implement me")
}

func (i *IdentityRepository) QueryEmail(context.Context, *identity.Email) error {
	// TODO implement me
	panic("implement me")
}
