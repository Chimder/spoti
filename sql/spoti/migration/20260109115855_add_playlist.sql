-- +goose Up
-- +goose StatementBegin

create table if not exists playlists (
    id uuid primary key default uuid_generate_v4 (),
    owner_id uuid not null references users (id) on delete cascade,
    playlist_name varchar(255) not null,
    description varchar(255),
    image varchar(255),
    is_public boolean not null default false,
    total int not null default 0,
    created_at timestamptz not null default now()
);

create index idx_playlists_owner_id on playlists (owner_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists playlists;

drop table if exists playlist_tracks;
-- +goose StatementEnd