-- noinspection SqlResolveForFile
SELECT
    active,
    COALESCE(authenticator_assurance_level, 0) AS authenticator_assurance_level,
    issued_at,
    expires_at,
    COALESCE(authenticated_at, 0::timestamptz) as authenticated_at,
    COALESCE(identity_id::text, '')  as identity_id
FROM sessions
WHERE id = $1;