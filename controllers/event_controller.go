package controllers

import (
	"net/http"

	"example-project.com/event-app/config"
	"example-project.com/event-app/models"
	"github.com/gin-gonic/gin"
)

func CreateEvent(ctx *gin.Context) {
	var event models.Event
	err := ctx.ShouldBindJSON(&event)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed create data",
			"error":   err.Error(),
		})
		return
	}

	// DUMMY
	event.UserID = 1

	config.DB.Create(&event)
	ctx.JSON(http.StatusCreated, gin.H{
		"data":    event,
		"message": "Successfully create data",
	})
}

func GetEvents(ctx *gin.Context) {
	var events []models.Event

	result := config.DB.Find(&events)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch events",
			"error":   result.Error.Error(),
		})
		return
	}

	// Jika aman, kirim response sukses
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Fetch all events successfully",
		"data":    events,
	})
}
