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

func GetEventById(ctx *gin.Context) {
	var event models.Event
	id := ctx.Param("id")

	var eventData = config.DB.First(&event, id).Error

	if eventData != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Event data details",
		"data":    event,
	})
}

func UpdateEvent(ctx *gin.Context) {
	var event models.Event

	id := ctx.Param("id")

	var eventData = config.DB.First(&event, id).Error

	if eventData != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	var input models.Event
	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed update data",
			"error":   err.Error(),
		})
		return
	}

	config.DB.Model(&event).Updates(input)
	ctx.JSON(http.StatusOK, gin.H{
		"data":    event,
		"message": "Successfully update data",
	})
}

func DeleteEvent(ctx *gin.Context) {
	var event models.Event
	id := ctx.Param("id")

	var eventData = config.DB.First(&event, id).Error
	if eventData != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	config.DB.Unscoped().Delete(&event)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully delete data",
		"data":    event,
	})
}
