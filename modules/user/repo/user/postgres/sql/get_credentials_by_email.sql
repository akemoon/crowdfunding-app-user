select c.user_id, c.password_hash, r.name
from credentials c
join roles r on r.id = c.role_id
where c.email = $1;
