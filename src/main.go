package main

import (
	"database/sql"

	"github.com/co-devs/ctf-api-emu/database"
	"github.com/co-devs/ctf-api-emu/handlers"
	"github.com/co-devs/ctf-api-emu/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	db = database.InitDB("ctf.db")
	defer db.Close()

	router := gin.Default()
	// router.Use(middleware.ApiKeyAuthMiddleware)

	router.GET("/heartbeat", handlers.GetHeartbeat)
	router.GET("/endpoints", middleware.ApiKeyAuthMiddleware, handlers.GetEndpoints)
	router.GET("/live_flags", middleware.ApiKeyAuthMiddleware, handlers.GetLiveFlags)
	router.GET("/submissions", middleware.ApiKeyAuthMiddleware, handlers.GetSubmittedFlags)
	router.POST("/submit", middleware.ApiKeyAuthMiddleware, handlers.PostFlag)
	router.GET("/status", middleware.ApiKeyAuthMiddleware, handlers.GetStatus)
	router.GET("/all_submissions", middleware.IsAdminAuthMiddleware, handlers.GetAllFlagSubmissions)
	router.GET("/secret", middleware.IsAdminAuthMiddleware, handlers.GetAPIKeys)

	router.Run("localhost:8080")
}
