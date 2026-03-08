-- +goose Up

create table if not exists credentials
(
    user_id       uuid primary key default uuidv7(),
    email         text not null,
    password_hash text not null,
    created_at    timestamptz not null default now(),
    -- TODO: is blocked

    constraint credentials_email_unique unique (email)
);

create table if not exists users
(
    id           uuid primary key,
    username     text not null,
    display_name text,
    description  text,
    created_at   timestamptz not null default now(),

    constraint users_username_unique unique (username)
);

create table if not exists follows
(
    follower_id uuid,
    followee_id uuid,
    created_at  timestamptz not null default now(),

    constraint follows_pk primary key (follower_id, followee_id),
    constraint follows_no_to_yourself check (follower_id <> followee_id),
    constraint follows_follower_fk foreign key (follower_id) references users(id),
    constraint follows_followee_fk foreign key (followee_id) references users(id)
);

-- +goose Down

drop table if exists follows;
drop table if exists users;
drop table if exists credentials;
