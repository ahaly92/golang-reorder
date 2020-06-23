package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/ahaly92/golang-reorder/drivers/sql"
	"github.com/ahaly92/golang-reorder/pkg/models"
	_ "github.com/lib/pq"
)

type postgresClient struct {
	pgxDriverWriter sql.Driver
	pgxDriverReader sql.Driver
}

func (pgClient postgresClient) GetAllUsers() (users []*models.User, err error) {
	rows, err := pgClient.pgxDriverReader.Query(context.Background(), getAllUsersQuery)

	if err != nil {
		return nil, err
	}
	if len(rows.Values) == 0 {
		return nil, nil
	}
	for _, row := range rows.Values {
		user := models.User{}
		err := pgClient.pgxDriverReader.Unmarshal(row,
			&user.ID,
			&user.Name,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (pgClient postgresClient) GetApplicationListForUser(userId int32) (applicationListItems []*models.ApplicationList, err error) {
	rows, err := pgClient.pgxDriverReader.Query(context.Background(), fmt.Sprintf(getApplicationListItemsForUser, userId))
	if err != nil {
		return nil, err
	}
	for _, row := range rows.Values {
		applicationListItem := models.ApplicationList{}
		err := pgClient.pgxDriverReader.Unmarshal(row,
			&applicationListItem.UserID,
			&applicationListItem.ApplicationID,
			&applicationListItem.Position,
		)
		if err != nil {
			return nil, err
		}

		applicationListItems = append(applicationListItems, &applicationListItem)
	}
	return applicationListItems, nil

}

func (pgClient postgresClient) AddUser(user models.User) error {
	_, err := pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(addUser, user.ID, user.Name))

	if err != nil {
		return err
	}
	return nil
}

func (pgClient postgresClient) AddApplication(description string) error {
	_, err := pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(addApplication, description))

	if err != nil {
		return err
	}
	return nil
}

func (pgClient postgresClient) DeleteApplication(applicationId int32) error {
	_, err := pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(deleteApplication, applicationId))

	if err != nil {
		return err
	}
	return nil
}

func (pgClient postgresClient) ReorderApplicationList(input models.ApplicationListInput) error {
	maxPositions, _ := pgClient.pgxDriverWriter.Query(context.Background(), fmt.Sprintf(getMaxItems, input.UserID))
	var maxPosition int32 = 0
	if len(maxPositions.Values) != 0 {
		_ = pgClient.pgxDriverReader.Unmarshal(maxPositions.Values[0],
			&maxPosition,
		)
	}

	rows, _ := pgClient.pgxDriverWriter.Query(context.Background(), fmt.Sprintf(getApplicationListItem, input.UserID, input.ApplicationID))
	if len(rows.Values) == 0 {
		_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(insertApplicationInList, input.UserID, input.ApplicationID, maxPosition+1))
		rows, _ = pgClient.pgxDriverWriter.Query(context.Background(), fmt.Sprintf(getApplicationListItem, input.UserID, input.ApplicationID))
		if input.DesiredPosition > maxPosition+1 {
			input.DesiredPosition = maxPosition + 1
		}
		if len(rows.Values) == 0 {
			return errors.New("unable to add application to list")
		}
	} else {
		if input.DesiredPosition > maxPosition {
			input.DesiredPosition = maxPosition
		}
	}

	applicationListItem := models.ApplicationList{}
	_ = pgClient.pgxDriverReader.Unmarshal(rows.Values[0],
		&applicationListItem.UserID,
		&applicationListItem.ApplicationID,
		&applicationListItem.Position,
	)

	if applicationListItem.Position != input.DesiredPosition {
		// move to position 0
		_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(setApplicationListItemPosition, 0, applicationListItem.Position, applicationListItem.UserID))

		// shift other items in list
		if input.DesiredPosition > applicationListItem.Position {
			print("down")
			_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(shiftApplicationListItemsDown, applicationListItem.Position, input.DesiredPosition))
		} else {
			print("up")
			_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(shiftApplicationListItemsUp, input.DesiredPosition, applicationListItem.Position))
		}

		// move to position to desired position
		_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(setApplicationListItemPosition, input.DesiredPosition, 0, applicationListItem.UserID))
	}

	return nil
}
