package cockroachdb

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.mreg.io/my-registry/auth/domain/identity"
	"gitlab.mreg.io/my-registry/auth/domain/session"
)

//go:embed sql/createSession.sql
var createSessionSQL string

//go:embed sql/querySessionByID.sql
var querySessionByIDSQL string

//go:embed sql/querySessionWithDevices.sql
var querySessionWithDevicesSQL string

//go:embed sql/updateDevice.sql
var updateDeviceSQL string

type sessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) session.Repository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *session.Session) error {
	if len(session.Devices) != 1 {
		return errors.New("only one device is allowed when creating a session")
	}

	device := session.Devices[0]

	return r.db.
		QueryRow(
			ctx,
			createSessionSQL,
			session.Active, zeronull.Int2(session.AuthenticatorAssuranceLevel), session.ExpiryInterval, zeronull.Timestamptz(session.AuthenticatedAt),
			session.Identity,
			device.IPAddress, device.GeoLocation, device.UserAgent,
		).
		Scan(&session.ID, &session.IssuedAt, &session.ExpiresAt, &session.Devices[0].ID)
}

func sessionFields(session *session.Session) []interface{} {
	return []interface{}{
		&session.Active,
		&session.AuthenticatorAssuranceLevel,
		&session.IssuedAt,
		&session.ExpiresAt,
		&session.AuthenticatedAt,
		&session.Identity.ID,
	}
}

func deviceFields(device *session.Device) []interface{} {
	return []interface{}{
		&device.ID,
		&device.IPAddress,
		&device.GeoLocation,
		&device.UserAgent,
		&device.SessionID,
	}
}

func (r *sessionRepository) QuerySessionByID(ctx context.Context, session *session.Session) error {
	if session.Identity == nil {
		session.Identity = &identity.Identity{}
	}
	return r.db.
		QueryRow(
			ctx,
			querySessionByIDSQL,
			session.ID,
		).
		Scan(sessionFields(session)...)
}

func (r *sessionRepository) QuerySessionWithDevices(ctx context.Context, sessionData *session.Session) error {
	if sessionData.Identity == nil {
		sessionData.Identity = &identity.Identity{}
	}
	err := r.db.
		QueryRow(
			ctx,
			querySessionByIDSQL,
			sessionData.ID,
		).
		Scan(sessionFields(sessionData)...)
	if err != nil {
		return err
	}
	rows, err := r.db.Query(ctx, querySessionWithDevicesSQL, sessionData.ID)
	if err != nil {
		return err
	}

	var devices []session.Device
	for rows.Next() {
		var device session.Device
		if err := rows.Scan(deviceFields(&device)...); err != nil {
			return err
		}
		device.SessionID = sessionData.ID
		devices = append(devices, device)
	}
	rows.Close()
	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return err
	}

	sessionData.Devices = devices
	return nil
}

func (r *sessionRepository) InsertDevice(ctx context.Context, device *session.Device) error {
	return r.db.
		QueryRow(
			ctx,
			updateDeviceSQL,
			device.IPAddress, device.GeoLocation, device.UserAgent, device.SessionID,
		).
		Scan(&device.ID)
}
