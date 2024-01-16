package handlers

import (
	"net/http"

	"spotify-api/db"
	"spotify-api/models"

	"github.com/gin-gonic/gin"
)

// HandleSearchByArtist retrieves tracks by artist
func HandleSearchByArtist(c *gin.Context) {
	var artist models.Artist
	artistName := c.Param("name")

	if err := db.DB.Where("name = ?", artistName).Preload("Tracks").First(&artist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}

	c.JSON(http.StatusOK, artist)
}
