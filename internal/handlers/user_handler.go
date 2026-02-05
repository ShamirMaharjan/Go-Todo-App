package handlers

import (
	"net/http"

	"github.com/ShamirMaharjan/Go-Todo-App/internal/models"
	"github.com/ShamirMaharjan/Go-Todo-App/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
			if err.Error() == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists" + err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": "User created successfully", "User": newUser})

	}
}
