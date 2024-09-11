package cockroachdb

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(ctx context.Context, connString string) *pgxpool.Pool {
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return conn
}
