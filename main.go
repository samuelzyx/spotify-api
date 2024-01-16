package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	clientID     = "6c84853b04074e9f8c8d2739c25a337a"
	clientSecret = "b6480532d600467797a17b0c16c16098"
	redirectURI  = "http://localhost:8080/callback"
	state        = "state-string"
	scopes       = []string{"user-read-email", "user-read-private"}
)

var oauthConfig = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  redirectURI,
	Scopes:       scopes,
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.spotify.com/authorize",
		TokenURL: "https://accounts.spotify.com/api/token",
	},
}

var authenticatedClient *http.Client
var db *gorm.DB

type Artist struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

type Album struct {
	gorm.Model
	ISRC       string `gorm:"uniqueIndex"`
	ImageURI   string
	Title      string
	ArtistName string
	Popularity int
}

func main() {
	//Connection DB
	var err error
	db, err = gorm.Open(mysql.Open("root:c8KidJWa@W&H59@tcp(localhost:3306)/spotify?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic("Error connecting to the database")
	}

	db.AutoMigrate(&Artist{}, &Album{})

	//Router Gin
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.GET("/", handleMain)
	r.GET("/login", handleLogin)
	r.GET("/callback", handleCallback)
	r.GET("/search/:isrc", handleSearchByISRC)
	r.GET("/search/artist/:name", handleSearchByArtistName)

	fmt.Println("Server is listening on :8080...")
	r.Run(":8080")
}

func handleMain(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func handleLogin(c *gin.Context) {
	url := oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleCallback(c *gin.Context) {
	code := c.Query("code")
	tok, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error exchanging code for access token")
		return
	}

	authenticatedClient = oauthConfig.Client(context.Background(), tok)
	c.String(http.StatusOK, "Successfully authenticated with Spotify. Access token: %s", tok.AccessToken)
}

func handleSearchByISRC(c *gin.Context) {
	isrc := c.Param("isrc")

	// Check if the song is already stored in the database
	// var existingAlbum Album
	// if err := db.Where("isrc = ?", isrc).First(&existingAlbum).Error; err == nil {
	// 	// The song already exists in the database, respond with the stored information
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"ISRC":       existingAlbum.ISRC,
	// 		"ImageURI":   existingAlbum.ImageURI,
	// 		"TitleAlbum": existingAlbum.Title,
	// 		// "ArtistID":   existingAlbum.ArtistID,
	// 		"ArtistName": existingAlbum.ArtistName,
	// 		"Popularity": existingAlbum.Popularity,
	// 	})
	// 	return
	// }

	// Create the Spotify search URL
	searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?type=track&q=isrc:%s", isrc)

	// Perform the request and get the response
	resp, err := authenticatedClient.Get(searchURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error making request to Spotify API")
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
		c.String(http.StatusInternalServerError, "Error getting the list of tracks from JSON response")
		return
	}

	// Extract the items from the list of tracks
	items, ok := tracks["items"].([]interface{})
	if !ok {
		c.String(http.StatusInternalServerError, "Error getting items from JSON response")
		return
	}

	// Initialize variables for the top artist and album with the highest popularity
	var topArtistName string
	var topTitle string
	var topImageURI string
	var topPopularity int

	// Iterate over the items to find the album with the highest popularity
	for _, item := range items {
		track, ok := item.(map[string]interface{})
		if !ok {
			c.String(http.StatusInternalServerError, "Error getting album information from JSON response")
			return
		}

		// Get the popularity of the album
		popularityFloat, ok := track["popularity"].(float64)
		if !ok {
			fmt.Println("Unable to get popularity, continue with the next album")
			continue // If unable to get popularity, continue with the next album
		}

		// If the popularity is higher than the current, update the variables
		popularity := int(popularityFloat)
		if popularity > topPopularity {
			topPopularity = popularity

			// Get information about the artist and album
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
		c.String(http.StatusInternalServerError, "Not track found, nothing stored")
		return
	}

	// Check if the artist is already stored in the database
	// var existingArtist Artist
	// if err := db.Where("name = ?", topArtistName).First(&existingArtist).Error; err != nil {
	// 	// The artist is not stored, store it
	// 	newArtist := Artist{
	// 		Name: topArtistName,
	// 	}
	// 	db.Create(&newArtist)

	// 	newAlbum := Album{
	// 		ISRC:       isrc,
	// 		ImageURI:   topImageURI,
	// 		Title:      topTitle,
	// 		ArtistID:   newArtist.ID,
	// 		Popularity: topPopularity,
	// 	}
	// 	db.Create(&newAlbum)

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"Message":    "New album added to the database",
	// 		"ISRC":       newAlbum.ISRC,
	// 		"ImageURI":   newAlbum.ImageURI,
	// 		"Title":      newAlbum.Title,
	// 		"ArtistName": newArtist.Name,
	// 		"Popularity": topPopularity,
	// 	})
	// 	return
	// }

	newAlbum := Album{
		ISRC:     isrc,
		ImageURI: topImageURI,
		Title:    topTitle,
		// ArtistID:   existingArtist.ID,
		ArtistName: topArtistName,
		Popularity: topPopularity,
	}
	db.Create(&newAlbum)

	c.JSON(http.StatusOK, gin.H{
		"Message":  "New album added to the database",
		"ISRC":     newAlbum.ISRC,
		"ImageURI": newAlbum.ImageURI,
		"Title":    newAlbum.Title,
		// "ArtistName": existingArtist.Name,
		"ArtistName": topArtistName,
		"Popularity": topPopularity,
	})
}

// Handler for searching albums by artist
func handleSearchByArtistName(c *gin.Context) {
	artistName := c.Param("artistName")

	// Perform the search in the database for albums associated with the artist
	var albums []Album
	if err := db.Where("name = ?", artistName).Find(&albums).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error searching for albums by artist")
		return
	}

	// Respond with the list of albums in JSON format
	c.JSON(http.StatusOK, albums)
}
