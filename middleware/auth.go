package middleware

import (
	"log"
	"net/http"
	"web-service-gin-tut/database"
	"web-service-gin-tut/models"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuthMiddleware validates API key from request header
func ApiKeyAuthMiddleware(c *gin.Context) {
	reqAPIKey := c.GetHeader("X-API-KEY")
	apiKey, err := getAPIKey(reqAPIKey)
	if err != nil || apiKey.Key == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid API key"})
		return
	}
	c.Set("apiKey", apiKey)
	c.Next()
}

// SecretAuthMiddleware checks if API key has permission to access the secret endpoint
func SecretAuthMiddleware(c *gin.Context) {
	apiKey, exists := c.Get("apiKey")
	if !exists || !apiKey.(models.APIKey).CanAccessSecret {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
		return
	}
	c.Next()
}

// PostAlbumAuthMiddleware checks if API key has permissions to create albums
func PostAlbumAuthMiddleware(c *gin.Context) {
	apiKey, exists := c.Get("apiKey")
	if !exists || !apiKey.(models.APIKey).CanAddAlbum {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
	}
	c.Next()
}

// ViewAlbumAuthMiddleware checks if the API key has permissions to view albums
func ViewAlbumAuthMiddleware(c *gin.Context) {
	apiKey, exists := c.Get("apiKey")
	if !exists || !apiKey.(models.APIKey).CanViewAlbum {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
	}
	c.Next()
}

// isValidAPIKey checks if the provided API key exists in the database
func isValidAPIKey(key string) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM api_keys WHERE key = ?)", key).Scan(&exists)
	if err != nil {
		log.Printf("Error checking API key: %v", err)
		return false
	}
	return exists
}

// getAPIKey retrieves the API key and its permissions from the database
func getAPIKey(key string) (models.APIKey, error) {
	var apiKey models.APIKey
	err := database.DB.QueryRow("SELECT key, can_view_secrets, can_add_album, can_view_album FROM api_keys WHERE key = ?", key).Scan(&apiKey.Key, &apiKey.CanAccessSecret, &apiKey.CanAddAlbum, &apiKey.CanViewAlbum)
	if err != nil {
		log.Printf("Error retrieving API key: %v", err)
		return apiKey, err
	}
	return apiKey, nil
}
