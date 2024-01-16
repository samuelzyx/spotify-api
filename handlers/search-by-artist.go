// Package handlers provides HTTP request handlers.
package handlers

import (
	"net/http"

	"spotify-api/db"
	"spotify-api/models"

	"github.com/gin-gonic/gin"
)

// @Summary Search tracks by artist
// @Description Get tracks by artist name
// @ID search-by-artist
// @Produce json
// @Param name path string true "Artist name"
// @Success 200 {object} models.Artist
// @Failure 404 {object} models.ErrorResponse
// @Router /search/{name} [get]
func HandleSearchByArtist(c *gin.Context) {
	var artist models.Artist
	artistName := c.Param("name")

	if err := db.DB.Where("name = ?", artistName).Preload("Tracks").First(&artist).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Code: http.StatusNotFound, Message: "Artist not found"})
		return
	}

	// Now artist.Tracks contains the associated tracks
	c.JSON(http.StatusOK, gin.H{
		"artist": artist,
		"tracks": artist.Tracks,
	})
}
