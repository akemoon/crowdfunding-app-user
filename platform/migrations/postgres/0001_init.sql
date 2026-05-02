-- +goose Up

create table if not exists roles
(
    id   smallint primary key,
    name text not null,

    constraint roles_name_unique unique (name)
);

insert into roles (id, name) values
(1, 'user'),
(2, 'moder'),
(3, 'admin');

create table if not exists credentials
(
    user_id       uuid primary key default uuidv7(),
    email         text not null,
    password_hash text not null,
    role_id       smallint not null default 1,
    created_at    timestamptz not null default now(),
    -- TODO: is blocked

    constraint credentials_email_unique unique (email),
    constraint credentials_role_fk foreign key (role_id) references roles(id)
);

create table if not exists default_avatars
(
    id  smallserial primary key,
    key text not null,

    constraint default_avatars_key_unique unique (key)
);

insert into default_avatars (key) values
    ('defaults/avatar-1.png'),
    ('defaults/avatar-2.png'),
    ('defaults/avatar-3.png'),
    ('defaults/avatar-4.png'),
    ('defaults/avatar-5.png'),
    ('defaults/avatar-6.png');

create table if not exists users
(
    id                uuid primary key,
    username          text not null,
    display_name      text,
    description       text,
    avatar_key        text,
    default_avatar_id smallint not null,
    created_at        timestamptz not null default now(),

    constraint users_username_unique unique (username),
    constraint users_default_avatar_fk foreign key (default_avatar_id) references default_avatars(id)
);

-- +goose StatementBegin
create or replace function assign_default_avatar()
returns trigger as $$
begin
    select id into new.default_avatar_id
    from default_avatars
    order by random()
    limit 1;

    return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

create trigger trigger_assign_default_avatar
    before insert on users
    for each row
    execute function assign_default_avatar();

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
drop trigger if exists trigger_assign_default_avatar on users;
drop function if exists assign_default_avatar;
drop table if exists users;
drop table if exists default_avatars;
drop table if exists credentials;
drop table if exists roles;
