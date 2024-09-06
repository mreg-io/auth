-- Handle inconsistency between CI and local migration during flyway clean
DROP TYPE IF EXISTS identity_state;
