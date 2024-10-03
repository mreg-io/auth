-- noinspection SqlResolveForFile
SELECT *
FROM devices
WHERE session_id = $1;