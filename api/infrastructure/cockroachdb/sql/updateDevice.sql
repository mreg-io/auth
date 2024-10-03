-- noinspection SqlResolveForFile
INSERT INTO devices (ip_address, geo_location, user_agent, session_id)
VALUES ($1, $2, $3, $4)
RETURNING id