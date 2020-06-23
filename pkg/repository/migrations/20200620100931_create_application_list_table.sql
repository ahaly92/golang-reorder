-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE application_lists (
    user_id int NOT NULL,
    application_id int NOT NULL,
    position int,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE application_lists;
-- +goose StatementEnd
