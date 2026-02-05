package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/ShamirMaharjan/Go-Todo-App/internal/config"
	"github.com/ShamirMaharjan/Go-Todo-App/internal/models"
	"github.com/ShamirMaharjan/Go-Todo-App/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var RegisterUser RegisterUserRequest

		var err error

		err = c.ShouldBindBodyWithJSON(&RegisterUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(RegisterUser.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be greater than 8 character"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RegisterUser.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to becrypt password" + err.Error()})
			return
		}

		user := &models.User{
			Email:    RegisterUser.Email,
			Password: string(hashedPassword),
		}

		newUser, err := repository.CreateUser(pool, user)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists \n " + err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": "User created successfully", "User": newUser})

	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest

		if err := c.ShouldBindBodyWithJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByEmail(pool, loginRequest.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentialss"})
			return
		}

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWT_SECRET))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})

	}
}
