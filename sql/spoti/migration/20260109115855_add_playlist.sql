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

create table if not exists user_saved_playlists (
    playlist_id uuid not null references playlists (id) on delete cascade,
    user_id uuid not null references users (id) on delete cascade,
    primary key (playlist_id, user_id)
);
create index idx_user_saved_playlist_user on user_saved_playlists (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_saved_playlists;
drop table if exists playlists;

-- +goose StatementEnd