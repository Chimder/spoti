-- +goose Up
-- +goose StatementBegin

create table if not exists albums (
    id uuid primary key default uuid_generate_v4 (),
    album_type varchar(225),
    total_tracks smallint not null,
    image varchar(255) not null,
    album_name varchar(255) not null unique,
    uri varchar(255) not null,
    copyrights varchar(255) not null,
    album_label varchar(255) not null,
    popularity smallint not null default 0,
    release_date timestamptz not null default now(),
    created_at timestamptz not null default now()
);

create table if not exists album_artists (
    album_id uuid not null references albums (id) on delete cascade,
    artist_id uuid not null references artists (id) on delete cascade,
    primary key (album_id, artist_id)
);

create index idx_album_artists_artist on album_artists (artist_id);

create table if not exists user_saved_albums (
    album_id uuid not null references albums (id) on delete cascade,
    user_id uuid not null references users (id) on delete cascade,
    primary key (album_id, user_id)
);

create index idx_user_saved_albums_user on user_saved_albums (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_saved_albums;

drop table if exists album_artists;

drop table if exists albums;
-- +goose StatementEnd