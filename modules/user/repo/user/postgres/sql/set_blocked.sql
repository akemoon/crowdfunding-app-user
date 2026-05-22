update credentials
set is_blocked = $2
where user_id = $1;
