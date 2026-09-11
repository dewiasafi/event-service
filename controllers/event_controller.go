package controllers

import (
	"context"
	"errors"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"example-project.com/event-app/config"
	"example-project.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	"gorm.io/gorm"
)

func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)
	return &client
}

func CreateEvent(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	// Menerima file form-data
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Image is required",
		})
		return
	}
	defer file.Close()

	// 1. upload file ke imageKit
	fileName := header.Filename
	ik := initImageKit()

	uploadResponse, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if errUpload != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Upload image failed",
		})
		return
	}

	parsedTime, _ := time.Parse(time.RFC3339, ctx.PostForm("datetime"))

	//Simpan ke database
	event := models.Event{
		Name:        ctx.PostForm("name"),
		Description: ctx.PostForm("description"),
		Location:    ctx.PostForm("location"),
		Datetime:    parsedTime,
		Image:       uploadResponse.URL,
		ImageID:     uploadResponse.FileID,
		UserID:      userId.(int),
	}

	config.DB.Create(&event)
	ctx.JSON(http.StatusCreated, gin.H{
		"data":    event,
		"message": "Successfully create data",
	})
}

func GetEvents(ctx *gin.Context) {
	// 1. Ambil userId dari context (diset oleh middleware auth/JWT)
	userId, exist := ctx.Get("userId")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var events []models.Event

	// 2. Inisiasi dasar GORM dengan filter berdasarkan user_id milik token
	query := config.DB.Model(&models.Event{}).Where("user_id = ?", userId)

	// 3. Tangkap fungsi filter by query (search)
	search := strings.TrimSpace(ctx.Query("search"))
	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchKeyword, searchKeyword) // Gunakan ILIKE agar tidak case-sensitive
	}

	// 4. Hitung total data (setelah difilter user_id & search, sebelum dilimit)
	var totalRow int64
	query.Count(&totalRow)

	// 5. Tangkap parameter pagination dan berikan nilai default
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, errPage := strconv.Atoi(pageStr)
	if errPage != nil || page < 1 {
		page = 1
	}

	limit, errLimit := strconv.Atoi(limitStr)
	if errLimit != nil || limit < 1 {
		limit = 10
	}

	// 6. Hitung offset dan total halaman
	offset := (page - 1) * limit
	totalPage := int(math.Ceil(float64(totalRow) / float64(limit)))

	// 7. Eksekusi query dengan Preload, Limit, dan Offset
	if err := query.Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fetch all events failed",
			"error":   err.Error(),
		})
		return
	}

	// 8. Kirim response sukses
	var resEvents []map[string]interface{}
	for _, e := range events {
		resEvents = append(resEvents, map[string]interface{}{
			"event_id":    e.ID,
			"name":        e.Name,
			"location":    e.Location,
			"datetime":    e.Datetime,
			"image":       e.Image,
			"description": e.Description,
			"user_id":     e.UserID,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Fetch all events successfully",
		"data": gin.H{
			"page":      page,
			"limit":     limit,
			"totalPage": totalPage,
			"total":     totalRow,
			"data":      resEvents,
		},
	})
}

func GetEventById(ctx *gin.Context) {
	// 1. Ambil userId dari context
	userId, exist := ctx.Get("userId")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var event models.Event
	id := ctx.Param("id")

	// 2. Eksekusi query dengan Preload User, cari berdasarkan ID dan pastikan milik user yang login
	err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Where("id = ? AND user_id = ?", id, userId).First(&event).Error

	// 3. Tangani jika error (data tidak ditemukan / error database)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Event not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch event detail",
			"error":   err.Error(),
		})
		return
	}

	resEvent := map[string]interface{}{
		"event_id":    event.ID,
		"name":        event.Name,
		"location":    event.Location,
		"datetime":    event.Datetime,
		"image":       event.Image,
		"description": event.Description,
		"user_id":     event.UserID,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Event data details",
		"data":    resEvent,
	})
}

func UpdateEvent(ctx *gin.Context) {
	var event models.Event

	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")

	var eventData = config.DB.First(&event, id).Error

	if eventData != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	if event.UserID != userId.(int) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User don't have access to update event",
		})
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		ik := initImageKit()

		fileName := header.Filename

		uploadResponse, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File:     file,
			FileName: fileName,
		})

		if errUpload == nil {
			// Hapus gambar lama
			if event.ImageID != "" {
				ik.Files.Delete(context.Background(), event.ImageID)
			}

			// Upload field gambar
			event.Image = uploadResponse.URL
			event.ImageID = uploadResponse.FileID
		}
	}

	if name := ctx.PostForm("name"); name != "" {
		event.Name = name
	}

	if description := ctx.PostForm("description"); description != "" {
		event.Description = description
	}

	if location := ctx.PostForm("location"); location != "" {
		event.Location = location
	}

	if dateTimeStr := ctx.PostForm("datetime"); dateTimeStr != "" {
		parseTime, errParse := time.Parse(time.RFC3339, dateTimeStr)
		if errParse == nil {
			event.Datetime = parseTime
		}
	}

	resEvent := map[string]interface{}{
		"event_id":    event.ID,
		"name":        event.Name,
		"location":    event.Location,
		"datetime":    event.Datetime,
		"image":       event.Image,
		"description": event.Description,
		"user_id":     event.UserID,
	}

	config.DB.Save(&event)
	ctx.JSON(http.StatusOK, gin.H{
		"data":    resEvent,
		"message": "Successfully update data",
	})
}

func DeleteEvent(ctx *gin.Context) {
	var event models.Event
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")

	var eventData = config.DB.First(&event, id).Error

	if eventData != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	if event.UserID != userId.(int) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User don't have access to delete event",
		})
		return
	}

	if event.ImageID != "" {
		ik := initImageKit()
		ik.Files.Delete(context.Background(), event.ImageID)
	}

	config.DB.Unscoped().Delete(&event)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully delete data",
	})
}
