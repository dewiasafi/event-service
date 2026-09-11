package controllers

import (
	"fmt"
	"net/http"
	"time"

	"example-project.com/event-app/config"
	"example-project.com/event-app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingInput struct {
	Phone   string `json:"phone" binding:"required"`
	EventID int    `json:"eventId" binding:"required"`
}

func CreateBookingEvent(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	var input BookingInput
	var booking models.Booking

	errValidation := ctx.ShouldBindJSON(&input)
	if errValidation != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	hasBooking := config.DB.Where("user_id = ? AND event_id = ?", userId.(int), input.EventID).First(&booking).Error
	if hasBooking == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "You have already booked this event", // Perbaiki typo juga
		})
		return
	}

	CodeBooking := fmt.Sprintf("BK-%sE%dU%d", time.Now().Format("20260101"), input.EventID, userId.(int))

	bookingData := models.Booking{
		Phone:       input.Phone,
		EventID:     input.EventID,
		BookingCode: CodeBooking,
		UserID:      userId.(int),
	}

	errCreateBooking := config.DB.Create(&bookingData).Error

	if errCreateBooking != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Booking failed",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Booking Success",
	})
}

func GetBookingByUser(ctx *gin.Context) {
	var booking []models.Booking
	userId, _ := ctx.Get("userId")

	errBook := config.DB.Preload("Event", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "description")
	}).Preload("Event.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Where("user_id = ?", userId).Find(&booking).Error

	if errBook != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch bookings",
		})
		return
	}

	var result []map[string]interface{}
	for _, b := range booking {
		result = append(result, map[string]interface{}{
			"booking_id":   b.ID,
			"phone":        b.Phone,
			"booking_code": b.BookingCode,
			"created_at":   b.CreatedAt,
			"event": map[string]interface{}{
				"name":        b.Event.Name,
				"description": b.Event.Description,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Fetch bookings success",
		"data":    result,
	})
}

func DeleteBooking(ctx *gin.Context) {
	var booking models.Booking
	userId, _ := ctx.Get("userId")
	paramsId := ctx.Param("id")

	bookingData := config.DB.First(&booking, paramsId).Error

	if bookingData != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	if booking.UserID != userId.(int) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You can't delete this book",
		})
		return
	}

	config.DB.Unscoped().Delete(&booking)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete booking success",
	})
}
