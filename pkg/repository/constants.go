package repository

const (
	getTodoListByUser = "SELECT * FROM " + todoListTableSchemaName + " WHERE userId ='%s'"

	todoListTableSchemaName  = "todoList"
	PostgresConnectionString = "host=%s port=%s user=%s password=%s dbname=%s"
)
