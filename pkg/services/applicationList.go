package services

import (
	"github.com/ahaly92/golang-reorder/pkg/models"
	"github.com/ahaly92/golang-reorder/pkg/repository"
)

type ApplicationListService interface {
	ReorderApplicationList(input models.ApplicationListInput) error
	GetApplicationListForUser(userId int32) (applicationListItems []*models.ApplicationList, err error)
	DeleteApplicationFromList(userId int32, applicationId int32) error
}

func NewApplicationListService(repo repository.Client) ApplicationListService {
	return &service{repo}
}

func (service *service) ReorderApplicationList(input models.ApplicationListInput) error {
	err := service.repo.ReorderApplicationList(input)
	if err != nil {
		return err
	}

	return nil
}

func (service *service) GetApplicationListForUser(userId int32) (applicationListItems []*models.ApplicationList, err error) {
	applicationListItems, err = service.repo.GetApplicationListForUser(userId)
	if err != nil {
		return nil, err
	}

	return applicationListItems, nil
}

func (service *service) DeleteApplicationFromList(userId int32, applicationId int32) error {
	err := service.repo.DeleteApplicationFromList(userId, applicationId)
	if err != nil {
		return err
	}

	return nil
}
