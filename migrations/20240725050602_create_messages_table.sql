-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id BIGSERIAL NOT NULL,
    message text,
    created_at timestamp,
    processed_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
