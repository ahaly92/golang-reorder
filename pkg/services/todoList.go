package services

import (
	"github.com/ahaly92/golang-reorder/pkg/models"
	"github.com/ahaly92/golang-reorder/pkg/repository"
)

type TodoListService interface {
	MoveTodoToInList(input models.TodoListInput) error
	GetTodoListForUser(userId int32) (todoListItems []*models.TodoList, err error)
}

func NewTodoListService(repo repository.Client) TodoListService {
	return &service{repo}
}

func (service *service) MoveTodoToInList(input models.TodoListInput) error {
	err := service.repo.MoveTodoToInList(input)
	if err != nil {
		return err
	}

	return nil
}

func (service *service) GetTodoListForUser(userId int32) (todoListItems []*models.TodoList, err error) {
	todoListItems, err = service.repo.GetTodoListForUser(userId)
	if err != nil {
		return nil, err
	}

	return todoListItems, nil
}
