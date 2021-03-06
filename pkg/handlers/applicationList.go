package handlers

import (
	"github.com/ahaly92/golang-reorder/pkg/models"
	"github.com/ahaly92/golang-reorder/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ReorderApplicationList(context *gin.Context, applicationListService services.ApplicationListService) {
	applicationListItem := models.ApplicationListInput{}
	_ = context.Bind(&applicationListItem)

	err := applicationListService.ReorderApplicationList(applicationListItem)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "application added / reordered to user's application list",
	})
}

func DeleteApplicationFromList(context *gin.Context, applicationListService services.ApplicationListService) {
	userId, _ := strconv.ParseInt(context.Param("userId"), 10, 32)
	applicationId, _ := strconv.ParseInt(context.Param("applicationId"), 10, 32)

	err := applicationListService.DeleteApplicationFromList(int32(userId), int32(applicationId))
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "application deleted from user's application list",
	})
}

func GetApplicationListForUser(context *gin.Context, applicationListService services.ApplicationListService) {
	userId, _ := strconv.ParseInt(context.Param("id"), 10, 32)
	applicationListItems, _ := applicationListService.GetApplicationListForUser(int32(userId))

	context.JSON(http.StatusOK, gin.H{
		"users": applicationListItems,
	})
}
