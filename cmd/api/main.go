package main

import (
	"log"

	"github.com/ShamirMaharjan/Go-Todo-App/internal/config"
	"github.com/ShamirMaharjan/Go-Todo-App/internal/database"
	"github.com/ShamirMaharjan/Go-Todo-App/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error

	cfg, err = config.Load()

	if err != nil {
		log.Printf("Unable to read .env file: %v", err)
	}

	var pool *pgxpool.Pool

	pool, err = database.Connect(cfg.DATABASE_URL)

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "TODO app is running!!",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/todo", handlers.CreateTodoHandler(pool))

	router.GET("/todo", handlers.GetAllTodosHandler(pool))

	router.GET("/todo/:id", handlers.GetTodoByIDHandler(pool))

	router.PUT("/todo/:id", handlers.UpdateTodoHandler(pool))

	router.DELETE("/todo/:id", handlers.DeleteTodoHandler(pool))

	router.POST("/auth/register", handlers.CreateUserHandler(pool))

	router.Run(":" + cfg.PORT)

}
