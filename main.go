package main

import (
	"database/sql"
	"log"
	"net/http"
	"web-service-gin-tut/database"
	"web-service-gin-tut/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)



var db *sql.DB

func main() {
	var err error
	db = database.InitDB("albums.db")
	defer db.Close()

	// create the api_keys table with permissions if it doesn't exist
	createAPIKeysTable := `CREATE TABLE IF NOT EXISTS api_keys (
		key TEXT PRIMARY KEY UNIQUE,
		can_view_secrets BOOLEAN NOT NULL DEFAULT 0
	);`
	if _, err = db.Exec(createAPIKeysTable); err != nil {
		log.Fatal(err)
	}

	// Insert API keys with permissions
	// insertAPIKey("test", false)
	// insertAPIKey("alfa", true)

	router := gin.Default()
	router.Use(apiKeyAuthMiddleware)

	router.GET("/albums", handlers.GetAlbums)
	router.GET("/albums/:id", handlers.GetAlbumByID)
	router.GET("/secret", secretAuthMiddleware, getAPIKeys)
	router.POST("/albums", handlers.PostAlbums)

	router.Run("localhost:8080")
}

// apiKeyAuthMiddleware validates API key from request header
func apiKeyAuthMiddleware(c *gin.Context) {
	reqAPIKey := c.GetHeader("X-API-KEY")
	if !isValidAPIKey(reqAPIKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid API key"})
		return
	}
	c.Set("apiKey", reqAPIKey)
	c.Next()
}

// secretAuthMiddleware checks if API key has permission to access the secret endpoint
func secretAuthMiddleware(c *gin.Context) {
	apiKey, exists := c.Get("apiKey")
	if !exists || !canAccessSecret(apiKey.(string)) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
		return
	}
	c.Next()
}

// isValidAPIKey checks if the provided API key exists in the database
func isValidAPIKey(key string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM api_keys WHERE key = ?)", key).Scan(&exists)
	if err != nil {
		log.Printf("Error checking API key: %v", err)
		return false
	}
	return exists
}

// canAccessSecret checks if the API key has permission to access the secret endpoint
func canAccessSecret(key string) bool {
	var canAccess bool
	err := db.QueryRow("SELECT can_view_secrets FROM api_keys WHERE key = ?", key).Scan(&canAccess)
	if err != nil {
		log.Printf("Error checking secret access: %v", err)
		return false
	}
	return canAccess
}

// insertAPIKey inserts an API key with permissions into the api_keys table if it doesn't exist
func insertAPIKey(key string, canAccessSecret bool) {
	_, err := db.Exec("INSERT OR IGNORE INTO api_keys (key, can_view_secrets) VALUES (?, ?)", key, canAccessSecret)
	if err != nil {
		log.Fatalf("Error inserting API key: %v", err)
	}
}

// getAPIKeys returns all API keys
func getAPIKeys(c *gin.Context) {
	rows, err := db.Query("SELECT key FROM api_keys;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var api_keys []string
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		api_keys = append(api_keys, a)
	}
	c.JSON(http.StatusOK, api_keys)
}
