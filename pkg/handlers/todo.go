package handlers

import (
	"github.com/ahaly92/golang-reorder/pkg/models"
	"github.com/ahaly92/golang-reorder/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddTodo(context *gin.Context, todoService services.TodoService) {
	todo := models.Todo{}
	_ = context.Bind(&todo)

	err := todoService.AddTodo(todo.Description)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "todo added!",
	})
}

func DeleteTodo(context *gin.Context, todoService services.TodoService) {
	todoId, _ := strconv.ParseInt(context.Param("id"), 10, 32)
	err := todoService.DeleteTodo(int32(todoId))
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "todo deleted!",
	})
}
