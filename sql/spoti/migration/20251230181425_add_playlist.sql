-- +goose Up
-- +goose StatementBegin

create table if not exists playlists (
    id uuid primary key default uuid_generate_v4(),
    owner_id uuid not null references users (id) on delete cascade,
    playlist_name varchar(255) not null,
    description varchar(255),
    disc_number smallint not null,
    image varchar(255),
    is_public boolean not null default false,
    total int not null default 0,
    created_at timestamptz not null default now()
);
create index idx_playlists_owner_id on playlists (owner_id);

create table if not exists playlist_tracks (
    playlist_id uuid not null references playlists (id) on delete cascade,
    track_id uuid not null references tracks (id) on delete cascade,
    track_position int not null,
    added_at timestamptz not null default now(),
    primary key (playlist_id, track_id),
    unique (playlist_id, track_position)
);

create index idx_playlist_tracks_playlist_id on playlist_tracks (playlist_id);
create index idx_playlist_tracks_track_id on playlist_tracks (track_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists playlists;
drop table if exists playlist_tracks;
-- +goose StatementEnd
