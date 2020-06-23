package handlers

import (
	"github.com/ahaly92/golang-reorder/pkg/models"
	"github.com/ahaly92/golang-reorder/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Users(context *gin.Context, userService services.UserService) {
	users, _ := userService.GetAllUsers()
	context.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func AddUser(context *gin.Context, userService services.UserService) {
	user := models.User{}
	_ = context.Bind(&user)

	err := userService.AddUser(user)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "user added!",
	})
}
