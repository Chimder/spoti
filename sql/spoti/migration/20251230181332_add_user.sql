-- +goose Up
-- +goose StatementBegin
create extension if not exists "uuid-ossp";

create table if not exists user (
    id uuid primary key default uuid_generate_v4(),
    user_name varchar(255) not null,
    email varchar(255) not null,
    image varchar(255) not null,
    followers bigint default 0,
    premium_status boolean default false

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists user;

-- +goose StatementEnd
