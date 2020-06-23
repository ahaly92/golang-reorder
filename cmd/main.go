package main

import (
	. "github.com/ahaly92/golang-reorder/pkg/handlers"
	"github.com/ahaly92/golang-reorder/pkg/repository"
	"github.com/ahaly92/golang-reorder/pkg/services"
	"github.com/gin-gonic/gin"
)

func main() {
	ginEngine := gin.Default()

	postgresClient, _ := repository.NewClient()

	userService := services.NewUserService(postgresClient)
	todoService := services.NewTodoService(postgresClient)
	todoListService := services.NewTodoListService(postgresClient)

	ginEngine.GET("/users", func(context *gin.Context) { Users(context, userService) })
	ginEngine.POST("/user", func(context *gin.Context) { AddUser(context, userService) })

	ginEngine.POST("/todo", func(context *gin.Context) { AddTodo(context, todoService) })
	ginEngine.DELETE("/todo/:id", func(context *gin.Context) { DeleteTodo(context, todoService) })

	ginEngine.POST("/todoList", func(context *gin.Context) { MoveTodoToInList(context, todoListService) })
	ginEngine.GET("/todoList/:id", func(context *gin.Context) { GetTodoListForUser(context, todoListService) })

	_ = ginEngine.Run(":4000")
}
