package handlers

import (
	"log"
	"net/http"
	"web-service-gin-tut/database"
	"web-service-gin-tut/models"

	"github.com/gin-gonic/gin"
)

// GetAPIKeys returns all API keys
func GetAPIKeys(c *gin.Context) {
	rows, err := database.DB.Query("SELECT key FROM api_keys;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var api_keys []models.APIKey
	for rows.Next() {
		var a models.APIKey
		if err := rows.Scan(&a.Key, &a.CanViewAlbum, &a.CanAddAlbum, &a.CanAccessSecret); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		api_keys = append(api_keys, a)
	}
	c.JSON(http.StatusOK, api_keys)
}

// InsertAPIKey inserts an API key with permissions into the api_keys table if it doesn't exist
func InsertAPIKey(key models.APIKey) {
	_, err := database.DB.Exec("INSERT OR IGNORE INTO api_keys (key, can_view_secrets, can_add_album, can_view_album) VALUES (?, ?, ?, ?)", key.Key, key.CanAccessSecret, key.CanAddAlbum, key.CanViewAlbum)
	if err != nil {
		log.Fatalf("Error inserting API key: %v", err)
	}
}

