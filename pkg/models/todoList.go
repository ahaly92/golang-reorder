package models

type TodoList struct {
	TodoID   int32 `json:"todo_id"`
	UserID   int32 `json:"user_id"`
	Position int32 `json:"position"`
}

type TodoListInput struct {
	TodoID          int32 `json:"todoId"`
	UserID          int32 `json:"userId"`
	Position        int32 `json:"position"`
	DesiredPosition int32 `json:"desiredPosition"`
}
