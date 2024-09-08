package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

const DatabaseUrlEnvName string = "DATABASE_URL"

var conn *pgx.Conn

func connectDatabase() {
	url, ok := os.LookupEnv(DatabaseUrlEnvName)
	if !ok {
		log.Fatalf("Cannot find %s environement variable\n", DatabaseUrlEnvName)
	}
	var err error
	conn, err = pgx.Connect(context.Background(), url) // TODO: os.Getenv("DATABASE_URL")
	if err != nil {
		panic("Unable to connect to database")
	}
	log.Println("Successfully Connected to DB.")
}
