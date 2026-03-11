-- +goose Up

-- Test users credentials:
-- alice.test@example.com / alice-test -> TestPass12345
-- boris.test@example.com / boris99    -> TestPass12345
-- cyril.test@example.com / cyril-dev  -> TestPass12345
-- strelok@zone.ua      / strelok     -> TestPass12345
insert into credentials (user_id, email, password_hash)
values
    ('00000000-0000-7000-8000-000000000001', 'alice.test@example.com', '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u'),
    ('00000000-0000-7000-8000-000000000002', 'boris.test@example.com', '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u'),
    ('00000000-0000-7000-8000-000000000003', 'cyril.test@example.com', '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u'),
    ('00000000-0000-7000-8000-000000000004', 'strelok@zone.ua',        '$2a$10$NiFFvrsWP0kiScpC1j3unO2iSdvva9BsqfV6JCP5dYyEyQwlnwb4u')
on conflict do nothing;

insert into users (id, username, display_name, description, avatar_key)
values
    ('00000000-0000-7000-8000-000000000001', 'alice-test', 'Alice Test', 'Test account for local development', ''),
    ('00000000-0000-7000-8000-000000000002', 'boris99', 'Борис Тест', 'Test account for local development', ''),
    ('00000000-0000-7000-8000-000000000003', 'cyril-dev', 'Cyril Dev', 'Test account for local development', ''),
    ('00000000-0000-7000-8000-000000000004', 'strelok', 'Стрелок', 'Меченый. Иду на Монолит.', '')
on conflict do nothing;

-- +goose Down

delete from users
where id in (
    '00000000-0000-7000-8000-000000000001',
    '00000000-0000-7000-8000-000000000002',
    '00000000-0000-7000-8000-000000000003',
    '00000000-0000-7000-8000-000000000004'
);

delete from credentials
where user_id in (
    '00000000-0000-7000-8000-000000000001',
    '00000000-0000-7000-8000-000000000002',
    '00000000-0000-7000-8000-000000000003',
    '00000000-0000-7000-8000-000000000004'
);
