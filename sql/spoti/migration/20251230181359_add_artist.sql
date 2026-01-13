-- +goose Up
-- +goose StatementBegin
create table if not exists artists (
    id uuid primary key default uuid_generate_v4(),
    url varchar(255) not null,
    uri varchar(255) not null,
    artist_name varchar(255) not null,
    image varchar(255),
    followers bigint default 0,
    popularity smallint default 0,
    genres text [] default '{}',
    artist_type varchar(50) not null,
    created_at timestamptz not null default now()
);
create index idx_artist_name on artists (artist_name);
create index idx_artist_popularity on artists (popularity);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists artists;

-- +goose StatementEnd
