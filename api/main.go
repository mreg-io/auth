package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	connectDatabase()
	defer func(conn *pgx.Conn, ctx context.Context) {
		if err := conn.Close(ctx); err != nil {
			log.Fatalf("Unable to close DB connection: %v", err)
		}
	}(conn, context.Background())

	bootstrap()

	mux := http.NewServeMux()
	registerHandlers(mux)

	// Start ConnectRPC server
	if err := http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
