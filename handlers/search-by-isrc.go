// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spotify-api/db"
	"spotify-api/models"

	"github.com/gin-gonic/gin"
)

// @Summary Search tracks by ISRC
// @Description Get track information by ISRC
// @ID search-by-isrc
// @Produce json
// @Param isrc path string true "ISRC code of the track"
// @Success 200 {object} models.Track
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /search/isrc/{isrc} [get]
func HandleSearchByISRC(c *gin.Context) {
	isrc := c.Param("isrc")

	// Check if the song is already stored in the database
	var existingTrack models.Track
	if err := db.DB.Preload("Artist").Where("isrc = ?", isrc).First(&existingTrack).Error; err == nil {
		// The song already exists in the database, respond with the stored information
		c.JSON(http.StatusOK, existingTrack)
		return
	}

	// Create the Spotify search URL
	searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?type=track&q=isrc:%s", isrc)

	// Perform the request and get the response
	resp, err := authenticatedClient.Get(searchURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusInternalServerError, Message: "Error making request to Spotify API"})
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding JSON response")
		return
	}

	// Extract the list of tracks from the JSON response
	tracks, ok := result["tracks"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusInternalServerError, Message: "Error getting the list of tracks from JSON response"})
		return
	}

	// Extract the items from the list of tracks
	items, ok := tracks["items"].([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusInternalServerError, Message: "Error getting tracks from JSON response"})
		return
	}

	// Initialize variables for the top artist and track with the highest popularity
	var topArtistName string
	var topTitle string
	var topImageURI string
	var topPopularity int

	// Iterate over the items to find the track with the highest popularity
	for _, item := range items {
		track, ok := item.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusInternalServerError, Message: "Error getting track information from JSON response"})
			return
		}

		// Get the popularity of the track
		popularityFloat, ok := track["popularity"].(float64)
		if !ok {
			fmt.Println("Unable to get popularity, continue with the next track")
			continue // If unable to get popularity, continue with the next track
		}

		// If the popularity is higher than the current, update the variables
		popularity := int(popularityFloat)
		if popularity > topPopularity {
			topPopularity = popularity

			// Get information about the artist and track
			artists, ok := track["artists"].([]interface{})
			if !ok || len(artists) == 0 {
				continue
			}
			artistInfo, ok := artists[0].(map[string]interface{})
			if !ok {
				continue
			}
			topArtistName, _ = artistInfo["name"].(string)

			topTitle, _ = track["name"].(string)

			images := track["album"].(map[string]interface{})["images"].([]interface{})
			imageURL := ""
			if len(images) > 0 {
				imageURL = images[0].(map[string]interface{})["url"].(string)
			}
			topImageURI = imageURL
		}
	}

	// Store the information in the database using GORM
	if topTitle == "" {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Code: http.StatusInternalServerError, Message: "Not track found, nothing stored"})
		return
	}

	var existingArtist models.Artist
	if err := db.DB.Preload("Tracks").Where("name = ?", topArtistName).First(&existingArtist).Error; err != nil {
		// The artist is not stored, store it
		newArtist := models.Artist{
			Name: topArtistName,
		}
		db.DB.Create(&newArtist)
		existingArtist = newArtist
	}

	newTrack := models.Track{
		ISRC:       isrc,
		ImageURI:   topImageURI,
		Title:      topTitle,
		Artist:     existingArtist, // Assign the artist to the track
		Popularity: topPopularity,
	}
	db.DB.Create(&newTrack)

	// Respond with the stored information
	c.JSON(http.StatusOK, existingTrack)
}
