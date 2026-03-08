select 
    id,
    username,
    display_name,
    description
from users
where id = $1;
