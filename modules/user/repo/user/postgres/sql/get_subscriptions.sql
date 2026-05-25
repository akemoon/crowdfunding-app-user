select
    u.id,
    u.username,
    u.display_name,
    u.description,
    case when u.avatar_key != '' then u.avatar_key else da.key end as avatar_key
from follows f
join users u on u.id = f.followee_id
join default_avatars da on da.id = u.default_avatar_id
where f.follower_id = $1
order by f.created_at desc;
