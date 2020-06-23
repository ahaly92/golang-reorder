package services

import "github.com/ahaly92/golang-reorder/pkg/repository"

type TodoService interface {
	AddTodo(description string) error
	DeleteTodo(todoId int32) error
}

func NewTodoService(repo repository.Client) TodoService {
	return &service{repo}
}

func (service *service) AddTodo(description string) error {
	err := service.repo.AddTodo(description)
	if err != nil {
		return err
	}

	return nil
}

func (service *service) DeleteTodo(todoId int32) error {
	err := service.repo.DeleteTodo(todoId)
	if err != nil {
		return err
	}

	return nil
}
