select 
    id,
    username,
    description
from users
where id = $1;
