select
    u.id,
    u.username,
    u.display_name,
    u.description,
    case when u.avatar_key != '' then u.avatar_key else da.key end as avatar_key,
    count(f.follower_id) as followers_count
from users u
join default_avatars da on da.id = u.default_avatar_id
left join follows f on f.followee_id = u.id
where u.id = $1
group by u.id, u.username, u.display_name, u.description, u.avatar_key, da.key;
