-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE applications (
    id SERIAL,
    description text,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE applications;
-- +goose StatementEnd
