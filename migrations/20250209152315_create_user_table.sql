-- +goose Up
create type role as enum ('USER', 'ADMIN');

create table users (
    id serial primary key,
    name text not null,
    email text not null,
    role role not null default 'USER',
    password text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
    );

-- +goose Down
drop table users;
