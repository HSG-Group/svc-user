CREATE TABLE users (
    id             UUID PRIMARY KEY,
    username       VARCHAR(32) NOT NULL UNIQUE,
    display_name   VARCHAR(80) NOT NULL,
    bio            VARCHAR(190) NOT NULL DEFAULT '',
    avatar_url     TEXT,
    custom_status  VARCHAR(128) NOT NULL DEFAULT '',
    discriminator  CHAR(4) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (username, discriminator)
);
