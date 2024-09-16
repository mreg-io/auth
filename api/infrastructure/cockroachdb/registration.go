package cockroachdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.mreg.io/my-registry/auth/domain/registration"
)

type RegistrationRepository struct {
	db *pgxpool.Pool
}

func NewRegistrationRepository(db *pgxpool.Pool) registration.Repository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) CreateFlow(ctx context.Context, flow *registration.Flow) error {
	// TODO
	panic("implement me")
}
