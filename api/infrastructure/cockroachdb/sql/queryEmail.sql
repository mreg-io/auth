-- noinspection SqlResolveForFile
SELECT verified, create_time, COALESCE(verified_at, 0::timestamptz), update_time
FROM emails
WHERE address = $1;