select c.user_id, c.password_hash, r.name, c.is_blocked
from credentials c
join roles r on r.id = c.role_id
where c.email = $1;
