-- +goose Up
-- +goose StatementBegin
INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
VALUES
('11111111-1111-1111-1111-111111111111', 'Netflix', 10, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2024-01-01', '2024-06-01'),
('22222222-2222-2222-2222-222222222222', 'Spotify', 5, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2024-02-01', '2024-08-01'),
('33333333-3333-3333-3333-333333333333', 'HBO', 12, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '2024-01-15', '2024-07-15'),
('44444444-4444-4444-4444-444444444444', 'Netflix', 10, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '2024-03-01', '2024-09-01'),
('55555555-5555-5555-5555-555555555555', 'Spotify', 5, 'cccccccc-cccc-cccc-cccc-cccccccccccc', '2024-04-01', NULL),
('66666666-6666-6666-6666-666666666666', 'HBO', 12, 'cccccccc-cccc-cccc-cccc-cccccccccccc', '2024-05-01', '2024-10-01'),
('77777777-7777-7777-7777-777777777777', 'Disney+', 8, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2024-01-10', '2024-06-10'),
('88888888-8888-8888-8888-888888888888', 'Amazon', 7, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '2024-02-10', '2024-07-10'),
('99999999-9999-9999-9999-999999999999', 'Netflix', 10, 'cccccccc-cccc-cccc-cccc-cccccccccccc', '2024-03-10', '2024-09-10'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab', 'Spotify', 5, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2024-04-15', '2024-10-15'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbc', 'HBO', 12, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2024-05-20', NULL),
('cccccccc-cccc-cccc-cccc-cccccccccccd', 'Disney+', 8, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '2024-06-01', '2024-12-01'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', 'Amazon', 7, 'cccccccc-cccc-cccc-cccc-cccccccccccc', '2024-07-01', '2024-12-31'),
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Netflix', 10, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2024-08-01', NULL),
('ffffffff-ffff-ffff-ffff-ffffffffffff', 'Spotify', 5, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '2024-09-01', '2025-03-01'),
('11111111-1111-1111-1111-111111111112', 'HBO', 12, 'cccccccc-cccc-cccc-cccc-cccccccccccc', '2024-10-01', '2025-04-01'),
('22222222-2222-2222-2222-222222222223', 'Disney+', 8, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2024-11-01', NULL),
('33333333-3333-3333-3333-333333333334', 'Amazon', 7, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '2024-12-01', NULL),
('44444444-4444-4444-4444-444444444445', 'Netflix', 10, 'cccccccc-cccc-cccc-cccc-cccccccccccc', '2025-01-01', '2025-06-01'),
('55555555-5555-5555-5555-555555555556', 'Spotify', 5, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2025-02-01', '2025-08-01');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM subscriptions
WHERE id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '44444444-4444-4444-4444-444444444444',
    '55555555-5555-5555-5555-555555555555',
    '66666666-6666-6666-6666-666666666666',
    '77777777-7777-7777-7777-777777777777',
    '88888888-8888-8888-8888-888888888888',
    '99999999-9999-9999-9999-999999999999',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbc',
    'cccccccc-cccc-cccc-cccc-cccccccccccd',
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
    'ffffffff-ffff-ffff-ffff-ffffffffffff',
    '11111111-1111-1111-1111-111111111112',
    '22222222-2222-2222-2222-222222222223',
    '33333333-3333-3333-3333-333333333334',
    '44444444-4444-4444-4444-444444444445',
    '55555555-5555-5555-5555-555555555556'
);
-- +goose StatementEnd
