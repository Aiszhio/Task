-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions (
                                             id            UUID        PRIMARY KEY,
                                             service_name  TEXT        NOT NULL,
                                             price         INTEGER     NOT NULL CHECK (price >= 0),
                                             user_id       UUID        NOT NULL,
                                             start_date    DATE        NOT NULL,
                                             end_date      DATE        NULL,
                                             created_at    TIMESTAMP   NOT NULL DEFAULT NOW(),
                                             updated_at    TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_service
    ON subscriptions (user_id, service_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscriptions_user_service;
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
