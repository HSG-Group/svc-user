CREATE TABLE user_preferences (
    user_id              UUID PRIMARY KEY REFERENCES users(id),
    theme                VARCHAR(10) NOT NULL DEFAULT 'dark',
    message_display      VARCHAR(10) NOT NULL DEFAULT 'cozy',
    notifications_dm     BOOL NOT NULL DEFAULT TRUE,
    notifications_mention BOOL NOT NULL DEFAULT TRUE,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);
