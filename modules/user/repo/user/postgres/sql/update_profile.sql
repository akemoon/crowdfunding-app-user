with updated as (
    update users
    set
        username = $2,
        display_name = $3,
        description = $4
    where id = $1
    returning id, username, display_name, description, avatar_key, default_avatar_id
)
select
    u.id,
    u.username,
    u.display_name,
    u.description,
    case when u.avatar_key != '' then u.avatar_key else da.key end as avatar_key
from updated u
join default_avatars da on da.id = u.default_avatar_id;
