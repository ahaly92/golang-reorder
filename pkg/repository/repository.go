package repository

import (
	"github.com/ahaly92/golang-reorder/drivers/sql"
	"github.com/ahaly92/golang-reorder/pkg/models"
)

type Client interface {
	GetAllUsers() (users []*models.User, err error)
	AddUser(user models.User) (err error)
	AddApplication(description string) error
	DeleteApplication(applicationId int32) error
	ReorderApplicationList(input models.ApplicationListInput) error
	GetApplicationListForUser(userId int32) (applicationListItems []*models.ApplicationList, err error)
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
