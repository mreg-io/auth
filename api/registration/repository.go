package registration

import "github.com/jackc/pgx/v5"

type Repository interface {
	insertFlow(flow *Flow) error
}

type repository struct {
}

func NewRepository(conn *pgx.Conn) Repository {
	return &repository{}
}

func (r *repository) insertFlow(flow *Flow) error {
	// TODO: sql operations
	// conn.Insert ...
	panic("implement me")
}

//TODO
