with updated as (
    update users
    set
        display_name = $2,
        description  = $3
    where id = $1
    returning id, username, display_name, description, avatar_key, default_avatar_id, updated_at
)
select
    u.id,
    u.username,
    coalesce(u.display_name, ''),
    coalesce(u.description, ''),
    case when u.avatar_key is not null and u.avatar_key != '' then u.avatar_key else da.key end as avatar_key,
    u.updated_at
from updated u
join default_avatars da on da.id = u.default_avatar_id;
