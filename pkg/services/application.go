package services

import "github.com/ahaly92/golang-reorder/pkg/repository"

type ApplicationService interface {
	AddApplication(description string) error
	DeleteApplication(applicationId int32) error
}

func NewApplicationService(repo repository.Client) ApplicationService {
	return &service{repo}
}

func (service *service) AddApplication(description string) error {
	err := service.repo.AddApplication(description)
	if err != nil {
		return err
	}

	return nil
}

func (service *service) DeleteApplication(applicationId int32) error {
	err := service.repo.DeleteApplication(applicationId)
	if err != nil {
		return err
	}

	return nil
}
