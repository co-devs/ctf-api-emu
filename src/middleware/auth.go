package middleware

import (
	"log"
	"net/http"

	"github.com/co-devs/ctf-api-emu/database"
	"github.com/co-devs/ctf-api-emu/models"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuthMiddleware validates API key from request header
func ApiKeyAuthMiddleware(c *gin.Context) {
	reqAPIKey := c.GetHeader("team-token")
	// log.Printf("API key: %v", reqAPIKey)
	apiKey, err := getAPIKey(reqAPIKey)
	if err != nil || apiKey.Key == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid API key"})
		return
	}
	c.Set("apiKey", apiKey)
	c.Next()
}

// IsAdminAuthMiddleware checks if the API key is an admin key
func IsAdminAuthMiddleware(c *gin.Context) {
	reqAPIKey := c.GetHeader("team-token")
	// log.Printf("API key: %v", reqAPIKey)
	apiKey, err := getAPIKey(reqAPIKey)
	if err != nil || !apiKey.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
	}
	c.Set("apiKey", apiKey)
	c.Next()
}

// isValidAPIKey checks if the provided API key exists in the database
func isValidAPIKey(key string) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM teams WHERE key = ?)", key).Scan(&exists)
	if err != nil {
		log.Printf("Error checking API key: %v", err)
		return false
	}
	return exists
}

// getAPIKey retrieves the API key and its permissions from the database
func getAPIKey(key string) (models.APIKey, error) {
	var apiKey models.APIKey
	// log.Printf("Query: SELECT id, key FROM teams WHERE key = '%v'", apiKey.Key)
	err := database.DB.QueryRow("SELECT id, key, is_admin FROM teams WHERE key = ?", key).Scan(&apiKey.TeamID, &apiKey.Key, &apiKey.IsAdmin)
	log.Printf("Resulting apiKey.TeamID %v, apiKey.Key %v, apiKey.IsAdmin %v", apiKey.TeamID, apiKey.Key, apiKey.IsAdmin)
	if err != nil {
		log.Printf("Error retrieving API key: %v", err)
		return apiKey, err
	}
	return apiKey, nil
}
