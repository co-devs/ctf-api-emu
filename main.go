package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// https://go.dev/doc/tutorial/web-service-gin
// album represents data about a record album
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func main() {
	var err error
	// initialize sqlite database connection
	db, err = sql.Open("sqlite3", "albums.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create the albums table if it doesn't exist
	createAlbumsTable := `CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		artist TEXT NOT NULL,
		price REAL NOT NULL
	);`
	if _, err = db.Exec(createAlbumsTable); err != nil {
		log.Fatal(err)
	}

	// create the api_keys table if it doesn't exist
	createAPIKeysTable := `CREATE TABLE IF NOT EXISTS api_keys (
		key TEXT PRIMARY KEY UNIQUE
	);`
	if _, err = db.Exec(createAPIKeysTable); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(apiKeyAuthMiddleware)

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

// apiKeyAuthMiddleware validates API key from request header
func apiKeyAuthMiddleware(c *gin.Context) {
	reqAPIKey := c.GetHeader("X-API-KEY")
	if !isValidAPIKey(reqAPIKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid API key"})
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

// insertAPIKey inserts an API key into the api_keys table if it doesn't exist
func insertAPIKey(key string) {
	_, err := db.Exec("INSERT OR IGNORE INTO api_keys (key) VALUES (?)", key)
	if err != nil {
		log.Fatalf("Error inserting API key: %v", err)
	}
}

// getAlbums response with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, artist, price FROM albums")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var albums []album
	for rows.Next() {
		var a album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		albums = append(albums, a)
	}
	c.JSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newAlbum.ID = strconv.FormatInt(id, 10)

	c.JSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var a album
	row := db.QueryRow("SELECT id, title, artist, price FROM albums WHERE id = ?", id)
	if err := row.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, a)
}
