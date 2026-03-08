delete from follows
where
    follower_id = $1
    and followee_id = $2;
