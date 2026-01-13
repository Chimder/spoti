-- +goose Up
-- +goose StatementBegin
create table tracks (
    id uuid primary key default uuid_generate_v4(),

    album_id uuid not null references albums (id) on delete cascade,
    recording_id uuid not null references recordings (id) on delete cascade,

    track_name varchar(255) not null,
    track_number smallint not null,
    disc_number smallint not null default 1,

    explicit boolean default false,
    is_playable boolean default true,

    track_type varchar(50) not null,        -- track / bonus / intro

    uri varchar(255) not null,
    islocal boolean default false,

    unique (album_id, disc_number, track_number)
);
create index idx_tracks_album_id on tracks (album_id);
create index idx_tracks_recording_id on tracks (recording_id);

create table if not exists artist_tracks (
    artist_id uuid not null references artists (id) on delete cascade,
    track_id uuid not null references tracks (id) on delete cascade,
    primary key (track_id, artist_id)
);

create index idx_track_artist_artist_id_hash on artist_tracks using hash (
    artist_id
);
create index idx_track_artist_track_id_hash on artist_tracks using hash (
    track_id
);

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
drop table if exists playlist_tracks;
drop table if exists artist_tracks;
drop table if exists tracks;
-- +goose StatementEnd
