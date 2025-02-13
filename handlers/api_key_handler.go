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
	rows, err := database.DB.Query("SELECT id, key, is_admin FROM teams;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var api_keys []models.APIKey
	for rows.Next() {
		var a models.APIKey
		if err := rows.Scan(&a.TeamID, &a.Key, &a.IsAdmin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		api_keys = append(api_keys, a)
	}
	c.JSON(http.StatusOK, api_keys)
}

// InsertAPIKey inserts an API key with permissions into the api_keys table if it doesn't exist
func InsertAPIKey(name string, key models.APIKey) {
	_, err := database.DB.Exec("INSERT OR IGNORE INTO api_keys (name, key, is_admin) VALUES (? ,?, ?)", name, key.Key, key.IsAdmin)
	if err != nil {
		log.Fatalf("Error inserting API key: %v", err)
	}
}

