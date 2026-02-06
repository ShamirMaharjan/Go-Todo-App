package handlers

import (
	"net/http"
	"strconv"

	"github.com/ShamirMaharjan/Go-Todo-App/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
		}

		userID := userIDInterface.(string)

		var input CreateTodoInput

		var err error

		err = c.ShouldBindBodyWithJSON(&input)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.CreateTodo(pool, input.Title, input.Completed, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, todo)

	}
}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
		}

		userID := userIDInterface.(string)

		todos, err := repository.GetAllTodos(pool, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}

func GetTodoByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
		}

		userID := userIDInterface.(string)

		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.GetTodoByID(pool, id, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)

	}
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
		}

		userID := userIDInterface.(string)

		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var input UpdateTodoInput

		err = c.ShouldBindBodyWithJSON(&input)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == nil && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least send title or completed"})
			return
		}

		exisiting, err := repository.GetTodoByID(pool, id, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		title := exisiting.Title

		if input.Title != nil {
			title = *input.Title
		}

		completed := exisiting.Completed

		if input.Completed != nil {
			completed = *input.Completed
		}

		todo, err := repository.UpdateTodoByID(pool, id, title, completed, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)

	}
}

func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
		}

		userID := userIDInterface.(string)

		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = repository.DeleteTodoByID(pool, id, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Success": "Todo Deleted Successfully"})
	}
}
