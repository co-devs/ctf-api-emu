package main

import (
	"database/sql"
	"web-service-gin-tut/database"
	"web-service-gin-tut/handlers"
	"web-service-gin-tut/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)



var db *sql.DB

func main() {
	db = database.InitDB("albums.db")
	defer db.Close()

	router := gin.Default()
	router.Use(middleware.ApiKeyAuthMiddleware)

	router.GET("/albums", handlers.GetAlbums)
	router.GET("/albums/:id", handlers.GetAlbumByID)
	router.GET("/secret", middleware.SecretAuthMiddleware, handlers.GetAPIKeys)
	router.POST("/albums", handlers.PostAlbums)

	router.Run("localhost:8080")
}
