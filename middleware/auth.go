package middleware

import (
	"log"
	"net/http"
	"web-service-gin-tut/database"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuthMiddleware validates API key from request header
func ApiKeyAuthMiddleware(c *gin.Context) {
	reqAPIKey := c.GetHeader("X-API-KEY")
	if !isValidAPIKey(reqAPIKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid API key"})
		return
	}
	c.Set("apiKey", reqAPIKey)
	c.Next()
}

// SecretAuthMiddleware checks if API key has permission to access the secret endpoint
func SecretAuthMiddleware(c *gin.Context) {
	apiKey, exists := c.Get("apiKey")
	if !exists || !canAccessSecret(apiKey.(string)) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
		return
	}
	c.Next()
}

// PostAlbumAuthMiddleware checks if API key has permissions to create albums
func PostAlbumAuthMiddleware(c *gin.Context) {
	apiKey, exists := c.Get("apiKey")
	if !exists || !canCreateAlbum(apiKey.(string)) {
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

// canAccessSecret checks if the API key has permission to access the secret endpoint
func canAccessSecret(key string) bool {
	var canAccess bool
	err := database.DB.QueryRow("SELECT can_view_secrets FROM api_keys WHERE key = ?", key).Scan(&canAccess)
	if err != nil {
		log.Printf("Error checking secret access: %v", err)
		return false
	}
	return canAccess
}

// canCreateAlbum checks if the API key has permission to create new albums
func canCreateAlbum (key string) bool {
	var canCreate bool
	err := database.DB.QueryRow("SELECT can_add_album FROM api_keys WHERE key = ?", key).Scan(&canCreate)
	if err != nil {
		log.Printf("Error checking add album access: %v", err)
		return false
	}
	return canCreate
}
