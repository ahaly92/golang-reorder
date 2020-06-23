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

func (pgClient postgresClient) GetTodoListForUser(userId int32) (todoListItems []*models.TodoList, err error) {
	rows, err := pgClient.pgxDriverReader.Query(context.Background(), fmt.Sprintf(getTodoListItemsForUser, userId))
	if err != nil {
		return nil, err
	}
	for _, row := range rows.Values {
		todoListItem := models.TodoList{}
		err := pgClient.pgxDriverReader.Unmarshal(row,
			&todoListItem.UserID,
			&todoListItem.TodoID,
			&todoListItem.Position,
		)
		if err != nil {
			return nil, err
		}

		todoListItems = append(todoListItems, &todoListItem)
	}
	return todoListItems, nil

}

func (pgClient postgresClient) AddUser(user models.User) error {
	_, err := pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(addUser, user.ID, user.Name))

	if err != nil {
		return err
	}
	return nil
}

func (pgClient postgresClient) AddTodo(description string) error {
	_, err := pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(addTodo, description))

	if err != nil {
		return err
	}
	return nil
}

func (pgClient postgresClient) DeleteTodo(todoId int32) error {
	_, err := pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(deleteTodo, todoId))

	if err != nil {
		return err
	}
	return nil
}

func (pgClient postgresClient) MoveTodoToInList(
	input models.TodoListInput) error {

	rows, _ := pgClient.pgxDriverWriter.Query(context.Background(), fmt.Sprintf(getTodoListItem, input.UserID, input.TodoID))

	if len(rows.Values) == 0 {
		return errors.New("todo list item not found")
	}

	todoListItem := models.TodoList{}
	_ = pgClient.pgxDriverReader.Unmarshal(rows.Values[0],
		&todoListItem.UserID,
		&todoListItem.TodoID,
		&todoListItem.Position,
	)

	// move to position 0
	_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(setTodoListItemPosition, 0, todoListItem.Position))

	// shift other items in list
	if input.DesiredPosition > todoListItem.Position {
		_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(shiftTodoListItemsDown, todoListItem.Position, input.DesiredPosition))
	} else {
		_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(shiftTodoListItemsUp, input.DesiredPosition, todoListItem.Position))
	}

	// move to position to desired position
	_, _ = pgClient.pgxDriverWriter.Exec(context.Background(), fmt.Sprintf(setTodoListItemPosition, input.DesiredPosition, 0))

	fmt.Printf("%+v\n", todoListItem)
	return nil
}
