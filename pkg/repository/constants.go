package repository

const (
	getAllUsersQuery = "SELECT * FROM " + usersTableName

	addUser = "INSERT INTO " + usersTableName + "(id, name) VALUES('%d', '%s')"

	addTodo    = "INSERT INTO " + todosTableName + "(description) VALUES('%s')"
	deleteTodo = "DELETE FROM " + todosTableName + " WHERE id='%d'"

	getTodoListItemsForUser = "SELECT * FROM " + todoListTableName + " WHERE user_id='%d'"

	getTodoListItem         = "SELECT * FROM " + todoListTableName + " WHERE user_id='%d' and todo_id='%d'"
	setTodoListItemPosition = "UPDATE " + todoListTableName + " SET position = '%d' WHERE position = '%d' AND user_id = 0;"
	shiftTodoListItemsDown  = "UPDATE todo_lists SET position = (position - 1) WHERE position > '%d' AND position <= '%d' AND user_id = user_id;"
	shiftTodoListItemsUp    = "UPDATE todo_lists SET position = (position + 1) WHERE position >= '%d' AND position < '%d' AND user_id = user_id;"

	usersTableName    = "users"
	todosTableName    = "todos"
	todoListTableName = "todo_lists"
)
