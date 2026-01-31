insert into users (id, username, description)
values ($1, $2, '')
returning id;
