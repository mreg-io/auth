-- noinspection SqlResolveForFile
SELECT issued_at, expires_at, session_id
FROM registration_flows
WHERE id = $1;