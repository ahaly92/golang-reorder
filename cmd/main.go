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
	applicationService := services.NewApplicationService(postgresClient)
	applicationListService := services.NewApplicationListService(postgresClient)

	ginEngine.GET("/users", func(context *gin.Context) { Users(context, userService) })
	ginEngine.POST("/user", func(context *gin.Context) { AddUser(context, userService) })

	ginEngine.POST("/application", func(context *gin.Context) { AddApplication(context, applicationService) })
	ginEngine.DELETE("/application/:id", func(context *gin.Context) { DeleteApplication(context, applicationService) })

	ginEngine.POST("/applicationList", func(context *gin.Context) { ReorderApplicationList(context, applicationListService) })
	ginEngine.GET("/applicationList/:id", func(context *gin.Context) { GetApplicationListForUser(context, applicationListService) })

	_ = ginEngine.Run(":4000")
}
