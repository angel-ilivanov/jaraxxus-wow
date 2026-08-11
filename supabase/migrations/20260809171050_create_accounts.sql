CREATE TABLE public.account (
    id BIGINT GENERATED ALWAYS AS IDENTITY (
        MINVALUE 1
        MAXVALUE 4294967295
    ) PRIMARY KEY,

    username TEXT NOT NULL,
    salt BYTEA NOT NULL,
    verifier BYTEA NOT NULL,
    session_key BYTEA,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT account_id_azeroth_range
        CHECK (id BETWEEN 1 AND 4294967295),

    CONSTRAINT account_username_unique
        UNIQUE (username),

    CONSTRAINT account_username_length
        CHECK (length(username) BETWEEN 1 AND 20),

    CONSTRAINT account_username_uppercase
        CHECK (username = upper(username)),

    CONSTRAINT account_salt_length
        CHECK (octet_length(salt) = 32),

    CONSTRAINT account_verifier_length
        CHECK (octet_length(verifier) = 32),

    CONSTRAINT account_session_key_length
        CHECK (octet_length(session_key) = 40)
);

COMMENT ON COLUMN public.account.updated_at IS
    'Maintained by application update statements using CURRENT_TIMESTAMP.';

ALTER TABLE public.account ENABLE ROW LEVEL SECURITY;

REVOKE ALL PRIVILEGES ON TABLE public.account
FROM anon, authenticated, service_role;
