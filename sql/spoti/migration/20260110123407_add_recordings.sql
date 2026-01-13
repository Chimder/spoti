-- +goose Up
-- +goose StatementBegin
create table recordings (
    id uuid primary key default uuid_generate_v4(),
    isrc varchar(15) not null unique,

    duration_ms bigint not null,

    popularity int not null default 0,
    play_count bigint not null default 0,

    audio_uri varchar(255) not null,
    preview_uri varchar(255),

    created_at timestamptz not null default now()
);
create index idx_recordings_id on recordings (id);
create index idx_recordings_isrc on recordings (isrc);
create index idx_recordings_popularity on recordings (popularity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists recordings;
-- +goose StatementEnd
