ALTER TABLE registration_flows
    DROP COLUMN identity_id,
    ADD COLUMN session_id UUID REFERENCES sessions (id) ON DELETE CASCADE;
CREATE INDEX ON registration_flows (session_id);
ALTER TABLE sessions
    ALTER COLUMN identity_id DROP NOT NULL;
ALTER TABLE devices ALTER COLUMN id SET DEFAULT gen_random_ulid();