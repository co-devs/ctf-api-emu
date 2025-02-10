package handlers

import (
	"log"
	"net/http"
	"web-service-gin-tut/database"

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

// InsertAPIKey inserts an API key with permissions into the api_keys table if it doesn't exist
func InsertAPIKey(key string, canAccessSecret bool) {
	_, err := database.DB.Exec("INSERT OR IGNORE INTO api_keys (key, can_view_secrets) VALUES (?, ?)", key, canAccessSecret)
	if err != nil {
		log.Fatalf("Error inserting API key: %v", err)
	}
}

