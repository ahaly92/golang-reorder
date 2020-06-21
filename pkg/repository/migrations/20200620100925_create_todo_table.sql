-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE todos (
    id int NOT NULL,
    description text,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE todos;
-- +goose StatementEnd
