-- +goose Up

create table if not exists users
(
    id          uuid primary key,
    username    text not null,
    description text,

    constraint users_username_unique unique (username)
);

create table if not exists subscriptions
(
    subscriber_id    uuid not null references users(id),
    subscribed_to_id uuid not null references users(id),

    constraint subscriptions_pk primary key (subscriber_id, subscribed_to_id),

    constraint subscriptions_no_self_subscribe check (subscriber_id <> subscribed_to_id)
);

-- +goose Down

drop table if exists subscriptions;
drop table if exists users;
