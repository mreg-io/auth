-- noinspection SqlResolveForFile
WITH inserted_flow AS (
    INSERT INTO registration_flows (expires_at, session_id)
        VALUES (current_timestamp + $1, $2)
        RETURNING id, issued_at, expires_at
)
SELECT id, issued_at, expires_at
FROM inserted_flow;
