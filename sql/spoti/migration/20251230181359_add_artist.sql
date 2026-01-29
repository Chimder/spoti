-- +goose Up
-- +goose StatementBegin
create table if not exists artists (
    id uuid primary key default uuid_generate_v4 (),
    url varchar(255) not null,
    uri varchar(255) not null,
    artist_name varchar(255) not null,
    image varchar(255),
    followers bigint default 0,
    popularity smallint default 0,
    genres text [] default '{}',
    created_at timestamptz not null default now()
);

create index idx_artist_name on artists (artist_name);
create index idx_artist_popularity on artists (popularity);


create table if not exists user_saved_artists (
    artist_id uuid not null references artists (id) on delete cascade,
    user_id uuid not null references users (id) on delete cascade,
    primary key (artist_id, user_id)
);
create index idx_user_saved_artists_user on user_saved_artists (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_saved_artists;
drop table if exists artists;

-- +goose StatementEnd