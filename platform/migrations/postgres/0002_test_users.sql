-- +goose Up

-- Test users credentials:
-- alice.test@example.com / alice-test -> TestPass12345
-- boris.test@example.com / boris99    -> TestPass12345
-- cyril.test@example.com / cyril-dev  -> TestPass12345
-- strelok@zone.ua        / strelok    -> TestPass12345
-- moder@example.com      / moder      -> TestPass12345
-- admin@example.com      / admin      -> TestPass12345
insert into users (id, username, display_name, description, avatar_key)
values
    ('00000000-0000-7000-8000-000000000001', 'alice-test', 'Alice Test', 'Test account for local development', ''),
    ('00000000-0000-7000-8000-000000000002', 'boris99', 'Борис Тест', 'Test account for local development', ''),
    ('00000000-0000-7000-8000-000000000003', 'cyril-dev', 'Cyril Dev', 'Test account for local development', ''),
    ('00000000-0000-7000-8000-000000000004', 'strelok', 'Стрелок', 'Меченый. Иду на Монолит.', 'defaults/avatar-1.png'),
    ('00000000-0000-7000-8000-000000000005', 'moder', 'Moder', 'Moderator account', ''),
    ('00000000-0000-7000-8000-000000000006', 'admin', 'Admin', 'Admin account', '')
on conflict do nothing;

insert into credentials (user_id, email, password_hash, role_id)
values
    ('00000000-0000-7000-8000-000000000001', 'alice.test@example.com', '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u', 1),
    ('00000000-0000-7000-8000-000000000002', 'boris.test@example.com', '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u', 1),
    ('00000000-0000-7000-8000-000000000003', 'cyril.test@example.com', '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u', 1),
    ('00000000-0000-7000-8000-000000000004', 'strelok@zone.ua',        '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u', 1),
    ('00000000-0000-7000-8000-000000000005', 'moder@example.com',      '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u', 2),
    ('00000000-0000-7000-8000-000000000006', 'admin@example.com',      '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u', 3)
on conflict do nothing;

-- +goose Down

delete from users
where id in (
    '00000000-0000-7000-8000-000000000001',
    '00000000-0000-7000-8000-000000000002',
    '00000000-0000-7000-8000-000000000003',
    '00000000-0000-7000-8000-000000000004',
    '00000000-0000-7000-8000-000000000005',
    '00000000-0000-7000-8000-000000000006'
);

delete from credentials
where user_id in (
    '00000000-0000-7000-8000-000000000001',
    '00000000-0000-7000-8000-000000000002',
    '00000000-0000-7000-8000-000000000003',
    '00000000-0000-7000-8000-000000000004',
    '00000000-0000-7000-8000-000000000005',
    '00000000-0000-7000-8000-000000000006'
);
