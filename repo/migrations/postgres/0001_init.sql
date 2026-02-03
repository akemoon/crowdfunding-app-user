-- +goose Up

create table if not exists users
(
    id          uuid primary key,
    username    text not null,
    description text,

    constraint users_username_unique unique (username)
);

-- +goose Down

drop table if exists users;
