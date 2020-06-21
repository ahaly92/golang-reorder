-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users (
    id int NOT NULL,
    first_name text,
    last_name text,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd
