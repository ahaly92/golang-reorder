-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE todo_lists (
    user_id int NOT NULL,
    todo_id int NOT NULL,
    position int,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE todo_lists;
-- +goose StatementEnd
