package registration

import "github.com/jackc/pgx/v5"

type Repository interface {
	insertFlow(flow *Flow) error
}

type repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) Repository {
	return &repository{conn}
}

func (r *repository) insertFlow(flow *Flow) error {
	// TODO: sql operations
	// conn.Insert ...
	panic("implement me")
}
