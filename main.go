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

	router.GET("/albums", middleware.ViewAlbumAuthMiddleware, handlers.GetAlbums)
	router.GET("/albums/:id", middleware.ViewAlbumAuthMiddleware, handlers.GetAlbumByID)
	router.GET("/secret", middleware.SecretAuthMiddleware, handlers.GetAPIKeys)
	router.POST("/albums", middleware.PostAlbumAuthMiddleware, handlers.PostAlbums)

	router.Run("localhost:8080")
}
