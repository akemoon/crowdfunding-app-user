update users
set
    username = $2,
    display_name = $3,
    description = $4
where id = $1
returning
    id,
    username,
    display_name,
    description;
