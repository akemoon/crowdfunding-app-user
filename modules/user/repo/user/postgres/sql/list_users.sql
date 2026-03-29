select
    u.id,
    u.username,
    coalesce(u.display_name, ''),
    coalesce(u.description, ''),
    case when u.avatar_key is not null and u.avatar_key != '' then u.avatar_key else da.key end as avatar_key,
    u.updated_at
from users u
join default_avatars da on da.id = u.default_avatar_id
order by u.id;
