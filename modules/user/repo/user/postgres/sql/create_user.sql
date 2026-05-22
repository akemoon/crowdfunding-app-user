insert into users (username, display_name, description, avatar_key)
values ($1, '', '', '')
returning id;
