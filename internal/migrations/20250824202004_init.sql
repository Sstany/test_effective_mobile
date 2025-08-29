-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    price int NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL
);

ALTER TABLE subscriptions
    ADD CONSTRAINT subscriptions_no_overlap
    EXCLUDE USING gist (
        user_id WITH =,
        title WITH =,
        daterange(start_date, COALESCE(end_date, 'infinity'::date), '[]') WITH &&
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
