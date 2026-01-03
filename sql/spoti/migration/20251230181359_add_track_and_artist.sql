-- +goose Up
-- +goose StatementBegin
create table if not exists artist (
    id uuid primary key default uuid_generate_v4(),
    url varchar(255) not null,
    uri varchar(255) not null,
    artist_name varchar(255) not null,
    image varchar(255),
    followers bigint default 0,
    popularity smallint default 0,
    genres text [] default '{}',
    artist_type varchar(50) not null
);
create index idx_artist_name on artist (artist_name);
create index idx_artist_popularity on artist (popularity);

create table if not exists track (
    id uuid primary key default uuid_generate_v4(),
    available_markets varchar(2) [] default '{}'::varchar(2) [],
    explicit boolean default false,
    is_playable boolean default false,
    track_name varchar(255) not null,
    popularity smallint not null,
    preview_url varchar(255) not null,
    disc_number smallint not null,
    track_number smallint not null,
    durationms bigint default 0,
    track_type varchar(255) not null,
    uri varchar(255) not null,
    islocal boolean default false
);

create table if not exists track_artist (
    artist_id uuid not null references artist (id) on delete cascade,
    track_id uuid not null references track (id) on delete cascade,
    primary key (track_id, artist_id)
);
create index idx_track_artist_artist_id_hash on track_artist using hash (
    artist_id
);
create index idx_track_artist_track_id_hash on track_artist using hash (
    track_id
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists track_artist;
drop table if exists track;
drop table if exists artist;

-- +goose StatementEnd
