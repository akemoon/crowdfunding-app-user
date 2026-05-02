select
    u.id,
    u.username,
    u.display_name,
    u.description,
    case when u.avatar_key != '' then u.avatar_key else da.key end as avatar_key
from users u
join default_avatars da on da.id = u.default_avatar_id
where ($1::text is null or u.username ilike '%' || $1 || '%')
order by u.username
limit $2 offset $3;
