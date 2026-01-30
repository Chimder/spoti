-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS spotify.listening_events (
    user_id String,
    track_id String,
    artist_id String,
    album_id String,
    duration_ms UInt32,
    is_skipped Bool,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree ()
PARTITION BY
    toYYYYMM (created_at)
ORDER BY (user_id, created_at, track_id) TTL created_at + INTERVAL 2 YEAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS spotify.listening_events
-- +goose StatementEnd