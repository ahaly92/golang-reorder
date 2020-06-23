package services

import (
	"github.com/ahaly92/golang-reorder/pkg/models"
	"github.com/ahaly92/golang-reorder/pkg/repository"
)

type service struct {
	repo repository.Client
}

type UserService interface {
	GetAllUsers() (users []*models.User, err error)
	AddUser(user models.User) error
}

func NewUserService(repo repository.Client) UserService {
	return &service{repo}
}

func (service *service) GetAllUsers() (users []*models.User, err error) {
	users, err = service.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *service) AddUser(user models.User) error {
	err := service.repo.AddUser(user)
	if err != nil {
		return err
	}

	return nil

}
