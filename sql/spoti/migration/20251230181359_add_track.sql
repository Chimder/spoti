-- +goose Up
-- +goose StatementBegin
create table if not exists track (
    id uuid primary key default uuid_generate_v4(),
    available_markets varchar(2) [] default array[]::varchar(2),
    explicit boolean default false,
    isplayable boolean default false,
    track_name varchar(255) not null,
    popularity smallint not null,
    previewurl varchar(255) not null,
    discnumber smallint not null,
    tracknumber smallint not null,
    durationms bigint default 0,
    track_type varchar(255) not null,
    uri varchar(255) not null,
    islocal boolean default false
);

create table if not exists track_artist (
    id uuid primary key default uuid_generate_v4(),
    artist_id uuid not null references artist (id) on delete cascade,
    track_id uuid not null references track (id) on delete cascade
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
drop table if exists track;
drop table if exists track_artist;

-- +goose StatementEnd
