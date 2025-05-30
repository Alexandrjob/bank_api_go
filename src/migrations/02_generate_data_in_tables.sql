-- 02_generate_data_in_tables.sql

-- +goose Up
TRUNCATE TABLE operations, users CASCADE;

INSERT INTO users (balance)
SELECT
    10000 * (random() * 2)::numeric(10,2)
FROM generate_series(1,50);

INSERT INTO operations (name, user_id, scope, date_create)
SELECT
    CASE WHEN random() < 0.4 THEN 'Пополнение'
         WHEN random() < 0.7 THEN 'Оплата услуг'
         ELSE 'Перевод' END,
    u.id,
    (1000 * (random() * 2))::numeric(10,2),
    NOW() - (random() * interval '365 days') + (random() * interval '23 hours')
FROM users u
    CROSS JOIN generate_series(1,20)
LIMIT 1000;