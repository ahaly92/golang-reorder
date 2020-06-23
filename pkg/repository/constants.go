package repository

const (
	getAllUsersQuery = "SELECT * FROM " + usersTableName

	addUser = "INSERT INTO " + usersTableName + "(id, name) VALUES('%d', '%s')"

	addApplication    = "INSERT INTO " + applicationsTableName + "(description) VALUES('%s')"
	deleteApplication = "DELETE FROM " + applicationsTableName + " WHERE id='%d'"

	getApplicationListItemsForUser = "SELECT * FROM " + applicationListTableName + " WHERE user_id='%d'"

	getApplicationListItem         = "SELECT * FROM " + applicationListTableName + " WHERE user_id='%d' and application_id='%d'"
	getMaxItems                    = "SELECT MAX(position) FROM " + applicationListTableName + " WHERE user_id='%d'"
	insertApplicationInList        = "INSERT INTO " + applicationListTableName + "(user_id, application_id, position) VALUES('%d', '%d', '%d')"
	setApplicationListItemPosition = "UPDATE " + applicationListTableName + " SET position = '%d' WHERE position = '%d' AND user_id = '%d';"
	shiftApplicationListItemsDown  = "UPDATE application_lists SET position = (position - 1) WHERE position > '%d' AND position <= '%d' AND user_id = user_id;"
	shiftApplicationListItemsUp    = "UPDATE application_lists SET position = (position + 1) WHERE position >= '%d' AND position < '%d' AND user_id = user_id;"

	usersTableName           = "users"
	applicationsTableName    = "applications"
	applicationListTableName = "application_lists"
)
