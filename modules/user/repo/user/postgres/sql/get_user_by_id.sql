select
    u.id,
    u.username,
    u.display_name,
    u.description,
    case when u.avatar_key != '' then u.avatar_key else da.key end as avatar_key
from users u
join default_avatars da on da.id = u.default_avatar_id
where u.id = $1;
