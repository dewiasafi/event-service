package main

import (
	"log"

	"fmt"

	"example-project.com/event-app/config"
	"example-project.com/event-app/controllers"
	"example-project.com/event-app/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	config.ConnectDB()

	server := gin.Default()
	server.Use(cors.Default())
	// Route
	api := server.Group("/api")
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.RegisterUser)
		auth.POST("/login", controllers.LoginUser)
	}

	protected := api.Group("/")
	protected.Use(middleware.RequiredAuth())
	{
		protected.GET("/auth/profil-user", controllers.GetUserLogin)
		protected.GET("/events", controllers.GetEvents)
		protected.POST("/event", controllers.CreateEvent)
		protected.GET("/event/:id", controllers.GetEventById)
		protected.PUT("/event/:id", controllers.UpdateEvent)
		protected.DELETE("/event/:id", controllers.DeleteEvent)

		protected.POST("/booking", controllers.CreateBookingEvent)
		protected.GET("/booking/user", controllers.GetBookingByUser)
		protected.DELETE("/booking/:id", controllers.DeleteBooking)

	}

	server.Run(":8080")
	fmt.Println("HELLO")
}
