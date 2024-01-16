package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"spotify-api/db"
	"spotify-api/handlers"
)

func main() {
	db.ConnectionDB()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", handlers.HandleWelcome)
	router.GET("/login", handlers.HandleLogin)
	router.GET("/callback", handlers.HandleCallback)
	router.GET("/search/:isrc", handlers.HandleSearchByISRC)
	router.GET("/search/artist/:name", handlers.HandleSearchByArtist)

	fmt.Println("Server is listening on :8080...")
	router.Run(":8080")
}
