package main

import (
	"log"

	"example-project.com/event-app/config"
	"example-project.com/event-app/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	server := gin.Default()

	// Route
	api := server.Group("/api")
	{
		api.GET("/events", controllers.GetEvents)
		api.POST("/event", controllers.CreateEvent)
	}

	server.Run(":8080")
}
