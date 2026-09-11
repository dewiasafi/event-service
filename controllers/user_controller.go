package controllers

import (
	"net/http"
	"os"
	"time"

	"example-project.com/event-app/config"
	"example-project.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func RegisterUser(ctx *gin.Context) {
	var input AuthInputRegister

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Encryption password failed",
		})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error

	if userCreated != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User Registration Failed",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User Registration Success",
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})
}

func LoginUser(ctx *gin.Context) {
	var input AuthInputLogin

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User

	userData := config.DB.Where("email = ?", input.Email).First(&user).Error
	if userData != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email belum terdaftar",
		})
		return
	}

	errMatchPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if errMatchPassword != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password wrong",
		})
		return
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 1).Unix(),
	})

	tokenString, errToken := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if errToken != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed build token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully Login",
		"token":   tokenString,
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})

}

func GetUserLogin(ctx *gin.Context) {
	userId, exits := ctx.Get("userId")
	if !exits {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var user models.User
	userData := config.DB.Select("id", "name", "email").First(&user, userId).Error
	if userData != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "fetch successfully",
		"data":    user,
	})
}
