package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Todos(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "todos coming soon",
	})
}
