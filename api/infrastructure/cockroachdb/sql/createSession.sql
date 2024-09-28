-- noinspection SqlResolveForFile
WITH session AS (
    INSERT INTO sessions (active, authenticator_assurance_level, expires_at, authenticated_at, identity_id)
        VALUES ($1, $2, current_timestamp + $3, $4, NULLIF($5::TEXT, '')::UUID)
        RETURNING *
),
device AS (
    INSERT INTO devices (ip_address, geo_location, user_agent, session_id)
        SELECT $6, $7, $8, session.id
        FROM session
    RETURNING *
)
SELECT session.id, issued_at, expires_at, device.id FROM session
    JOIN device ON session.id = device.session_id
;