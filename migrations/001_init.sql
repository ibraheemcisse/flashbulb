
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    handle TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT users_handle_not_empty CHECK (length(trim(handle)) > 0)
);

CREATE UNIQUE INDEX users_handle_unique
    ON users (lower(handle));


-- A drop defines the period during which users can publish.
CREATE TABLE drops (
    id BIGSERIAL PRIMARY KEY,
    opens_at TIMESTAMPTZ NOT NULL,
    closes_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT drops_valid_window CHECK (closes_at > opens_at)
);


-- Every post belongs to a drop and has its own expiry time.
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    drop_id BIGINT NOT NULL REFERENCES drops(id) ON DELETE RESTRICT,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT posts_valid_expiry CHECK (expires_at > created_at)
);

CREATE INDEX posts_feed_idx
    ON posts (created_at DESC, id DESC);

CREATE INDEX posts_user_created_idx
    ON posts (user_id, created_at DESC);

CREATE INDEX posts_expires_idx
    ON posts (expires_at);


CREATE TABLE likes (
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (post_id, user_id)
);

CREATE INDEX likes_user_idx
    ON likes (user_id, created_at DESC);


CREATE TABLE photos (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    object_key TEXT NOT NULL UNIQUE,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT messages_different_participants
        CHECK (sender_id <> recipient_id),

    CONSTRAINT messages_body_not_empty
        CHECK (length(trim(body)) > 0)
);

CREATE INDEX messages_recipient_idx
    ON messages (recipient_id, created_at DESC);

CREATE INDEX messages_sender_idx
    ON messages (sender_id, created_at DESC);
