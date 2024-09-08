package session

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	InsertSessionAndDevice(session *Session, device *Device) error
	AddDevice(device *Device) error
}

type repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) Repository {
	return &repository{conn}
}

func (r *repository) InsertSessionAndDevice(session *Session, device *Device) error {
	// TODO implement me
	err := r.conn.QueryRow(context.Background(),
		"INSERT INTO sessions(active, authenticator_assurance_level, issued_at, expires_at) VALUES ($1, $2, $3, $4) RETURNING id",
		session.Active, session.AuthenticatorAssuranceLevel, session.IssuedAt, session.ExpiresAt).Scan(&session.ID)
	if err != nil {
		log.Fatalf("Error inserting session in database: %v", err)
	}
	_, err = r.conn.Exec(context.Background(),
		"INSERT INTO devices (ip_address, geo_location, user_agent, session_id) VALUES ($1, $2, $3, $4)",
		device.IPAddress, device.GeoLocation, device.UserAgent, session.ID)
	if err != nil {
		log.Fatalf("Error adding device in database: %v", err)
	}

	return nil
}

func (r *repository) AddDevice(device *Device) error {
	_, err := r.conn.Exec(context.Background(),
		"INSERT INTO devices (ip_address, geo_location, user_agent, session_id) VALUES ($1, $2, $3, $4)",
		device.IPAddress, device.GeoLocation, device.UserAgent, device.SessionID)
	if err != nil {
		log.Fatalf("Error adding device in database: %v", err)
	}
	return nil
}
