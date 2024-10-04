-- noinspection SqlResolveForFile
WITH
    identity AS (
        INSERT INTO identities(timezone)
        VALUES ($1)
        RETURNING *
    ),
    email AS (
        INSERT INTO emails (address, identity_id)
        SELECT  $2, identity.id
        FROM identity
        RETURNING *
    ),
    password AS (
        INSERT INTO passwords (identity_id, password_hash)
        SELECT identity.id, $3
        FROM identity
        RETURNING *
    )
SELECT
    identity.id,
    identity.create_time,
    identity.update_time,
    identity.state_update_time,
    email.create_time,
    email.update_time
FROM identity, email
;
