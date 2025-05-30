-- 01_create_tables.sql

-- +goose Up
CREATE TABLE users
(
    id      BIGSERIAL PRIMARY KEY,
    balance DECIMAL NOT NULL DEFAULT 0.00
);

CREATE TABLE operations
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    user_id     BIGINT         NOT NULL,
    scope       DECIMAL NOT NULL,
    date_create TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE INDEX idx_operations_user_id ON operations (user_id);
CREATE INDEX idx_operations_date_create ON operations (date_create);

-- +goose Down
--DROP TABLE IF EXISTS operations;
--DROP TABLE IF EXISTS users;