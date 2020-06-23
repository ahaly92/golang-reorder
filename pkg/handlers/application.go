package handlers

import (
	"github.com/ahaly92/golang-reorder/pkg/models"
	"github.com/ahaly92/golang-reorder/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddApplication(context *gin.Context, applicationService services.ApplicationService) {
	application := models.Application{}
	_ = context.Bind(&application)

	err := applicationService.AddApplication(application.Description)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "application added!",
	})
}

func DeleteApplication(context *gin.Context, applicationService services.ApplicationService) {
	applicationId, _ := strconv.ParseInt(context.Param("id"), 10, 32)
	err := applicationService.DeleteApplication(int32(applicationId))
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "application deleted!",
	})
}
