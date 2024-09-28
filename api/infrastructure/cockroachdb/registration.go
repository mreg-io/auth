package cockroachdb

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.mreg.io/my-registry/auth/domain/registration"
)

//go:embed sql/createRegistrationFlow.sql
var insertRegistrationFlowSQL string

type RegistrationRepository struct {
	db *pgxpool.Pool
}

func NewRegistrationRepository(db *pgxpool.Pool) registration.Repository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) CreateFlow(ctx context.Context, flow *registration.Flow) error {
	return r.db.
		QueryRow(
			ctx,
			insertRegistrationFlowSQL,
			flow.Interval, flow.SessionID,
		).
		Scan(&flow.FlowID, &flow.IssuedAt, &flow.ExpiresAt)
}

func (r *RegistrationRepository) QueryFlow(context.Context, *registration.Flow) error {
	// TODO: Implement me
	panic("not implemented")
}
