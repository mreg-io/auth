CREATE TYPE auth.public.identity_state AS ENUM ('active', 'suspended');

CREATE TABLE IF NOT EXISTS auth.public.identities
(
    id                UUID PRIMARY KEY                    DEFAULT gen_random_ulid(),
    state             auth.public.identity_state NOT NULL DEFAULT 'active',
    full_name         STRING(64),
    display_name      STRING(64),
    avatar            STRING(256),
    timezone          STRING(64) NOT NULL,
    create_time       TIMESTAMPTZ                NOT NULL DEFAULT current_timestamp(),
    update_time       TIMESTAMPTZ                NOT NULL DEFAULT current_timestamp(),
    state_update_time TIMESTAMPTZ                NOT NULL DEFAULT current_timestamp()
);