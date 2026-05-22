select c.user_id, c.email, c.password_hash, r.name, c.is_blocked
from credentials c
join roles r on r.id = c.role_id
where c.user_id = $1;
