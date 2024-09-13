-- noinspection SqlResolveForFile
WITH session AS (
    INSERT INTO sessions (active, authenticator_assurance_level, expires_at, authenticated_at)
        VALUES ($1, $2, current_timestamp + $3, $4)
        RETURNING *
),
device AS (
    INSERT INTO devices (ip_address, geo_location, user_agent, session_id)
        SELECT $5, $6, $7, session.id
        FROM session
    RETURNING *
)
SELECT session.id, issued_at, expires_at, device.id FROM session
    JOIN device ON session.id = device.session_id
;