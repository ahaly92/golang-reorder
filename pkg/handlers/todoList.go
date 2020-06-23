package handlers

import (
	"github.com/ahaly92/golang-reorder/pkg/models"
	"github.com/ahaly92/golang-reorder/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func MoveTodoToInList(context *gin.Context, todoListService services.TodoListService) {
	todoListItem := models.TodoListInput{}
	_ = context.Bind(&todoListItem)

	err := todoListService.MoveTodoToInList(todoListItem)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "todo added to TodoList!",
	})
}

func GetTodoListForUser(context *gin.Context, todoListService services.TodoListService) {
	userId, _ := strconv.ParseInt(context.Param("id"), 10, 32)
	todoListItems, _ := todoListService.GetTodoListForUser(int32(userId))

	context.JSON(http.StatusOK, gin.H{
		"users": todoListItems,
	})
}
