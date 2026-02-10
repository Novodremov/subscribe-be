-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_name  TEXT        NOT NULL,
    price         INTEGER     NOT NULL CHECK (price > 0),
    user_id       UUID        NOT NULL,
    start_date    DATE        NOT NULL,
    end_date      DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (end_date IS NULL OR end_date >= start_date)
);

-- Индексы под основные сценарии
CREATE INDEX IF NOT EXISTS subscriptions_user_id_idx
    ON subscriptions (user_id);

CREATE INDEX IF NOT EXISTS subscriptions_service_name_idx
    ON subscriptions (service_name);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
