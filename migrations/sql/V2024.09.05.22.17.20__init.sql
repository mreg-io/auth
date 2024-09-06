CREATE TYPE identity_state AS ENUM ('active', 'suspended');

CREATE TABLE identities
(
    id                UUID PRIMARY KEY                                              DEFAULT gen_random_ulid(),
    state             identity_state                           NOT NULL DEFAULT 'active',
    full_name         STRING(64),
    display_name      STRING(64),
    avatar            STRING(256),
    timezone          STRING(64) NOT NULL,
    create_time       TIMESTAMPTZ                                          NOT NULL DEFAULT current_timestamp(),
    update_time       TIMESTAMPTZ CHECK (update_time >= create_time)       NOT NULL DEFAULT current_timestamp(),
    state_update_time TIMESTAMPTZ CHECK (state_update_time >= create_time) NOT NULL DEFAULT current_timestamp()
);

CREATE TABLE sessions
(
    id                            UUID PRIMARY KEY                                     DEFAULT gen_random_ulid(),
    active                        BOOL                                        NOT NULL DEFAULT true,
    authenticator_assurance_level SMALLINT CHECK (authenticator_assurance_level IS NULL OR
                                                  authenticator_assurance_level BETWEEN 1 AND 3),
    issued_at                     TIMESTAMPTZ                                 NOT NULL DEFAULT current_timestamp(),
    expires_at                    TIMESTAMPTZ CHECK (expires_at >= issued_at) NOT NULL,
    authenticated_at              TIMESTAMPTZ CHECK (authenticated_at >= issued_at AND authenticated_at <= expires_at),
    identity_id                   UUID                                        NOT NULL REFERENCES identities (id) ON DELETE CASCADE,
    INDEX identity_id_idx (identity_id)
);

CREATE TABLE emails
(
    address     STRING(320) PRIMARY KEY,
    verified    BOOLEAN                                        NOT NULL DEFAULT false,
    create_time TIMESTAMPTZ                                    NOT NULL DEFAULT current_timestamp(),
    verified_at TIMESTAMPTZ CHECK (verified_at >= create_time ),
    update_time TIMESTAMPTZ CHECK (update_time >= create_time) NOT NULL DEFAULT current_timestamp(),
    identity_id UUID                                           NOT NULL REFERENCES identities (id) ON DELETE CASCADE,
    INDEX identity_id_idx (identity_id)
);

CREATE TABLE phones
(
    number      STRING(16) PRIMARY KEY,
    verified    BOOLEAN                                        NOT NULL DEFAULT false,
    create_time TIMESTAMPTZ                                    NOT NULL DEFAULT current_timestamp(),
    verified_at TIMESTAMPTZ CHECK (verified_at >= create_time ),
    update_time TIMESTAMPTZ CHECK (update_time >= create_time) NOT NULL DEFAULT current_timestamp(),
    identity_id UUID                                           NOT NULL REFERENCES identities (id) ON DELETE CASCADE,
    INDEX identity_id_idx (identity_id)
);

CREATE TABLE registration_flows
(
    id          UUID PRIMARY KEY                                     DEFAULT gen_random_ulid(),
    issued_at   TIMESTAMPTZ                                 NOT NULL DEFAULT current_timestamp(),
    expires_at  TIMESTAMPTZ CHECK (expires_at >= issued_at) NOT NULL,
    identity_id UUID                                        NOT NULL REFERENCES identities (id) ON DELETE CASCADE,
    INDEX identity_id_idx (identity_id)
);

CREATE TABLE passwords
(
    identity_id   UUID PRIMARY KEY REFERENCES identities (id) ON DELETE CASCADE,
    password_hash STRING(256) NOT NULL
);

CREATE TABLE authentication_methods
(
    id          UUID PRIMARY KEY                   DEFAULT gen_random_ulid(),
    aal         SMALLINT CHECK (aal >= 1) NOT NULL,
    complete_at TIMESTAMPTZ               NOT NULL DEFAULT current_timestamp(),
    password_id UUID                      NOT NULL REFERENCES passwords (identity_id) ON DELETE CASCADE,
    session_id  UUID                      NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    INDEX session_id_idx (session_id)
);

CREATE TABLE devices
(
    id           UUID PRIMARY KEY,
    ip_address   INET NOT NULL,
    geo_location STRING(64) NOT NULL,
    user_agent   STRING(256) NOT NULL,
    session_id   UUID NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    INDEX session_id_idx (session_id)
);
