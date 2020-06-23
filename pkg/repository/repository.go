package repository

import (
	"github.com/ahaly92/golang-reorder/drivers/sql"
	"github.com/ahaly92/golang-reorder/pkg/models"
)

type Client interface {
	GetAllUsers() (users []*models.User, err error)
	AddUser(user models.User) (err error)
	AddTodo(description string) error
	DeleteTodo(todoId int32) error
	MoveTodoToInList(input models.TodoListInput) error
	GetTodoListForUser(userId int32) (todoListItems []*models.TodoList, err error)
}

func NewClient() (Client, error) {
	pgxDriver, err := sql.CreatePostgresConnection(
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"reorder",
		true,
		30,
		200)
	if err != nil {
		return nil, err
	}
	return &postgresClient{pgxDriverWriter: pgxDriver, pgxDriverReader: pgxDriver}, nil
}
