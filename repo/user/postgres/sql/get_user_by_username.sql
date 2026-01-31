select 
    id,
    username,
    description
from users
where username = $1;
