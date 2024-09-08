package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

const DatabaseURLEnvName string = "DATABASE_URL"

var conn *pgx.Conn

func connectDatabase() {
	url, ok := os.LookupEnv(DatabaseURLEnvName)
	if !ok {
		log.Fatalf("Cannot find %s environement variable\n", DatabaseURLEnvName)
	}
	var err error
	conn, err = pgx.Connect(context.Background(), url)
	if err != nil {
		panic("Unable to connect to database")
	}
	log.Println("Successfully Connected to DB.")
}
