package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"spotify-api/controllers"
	"spotify-api/db"
)

func main() {
	//Connection DB
	db.ConnectionDB()

	//Router Gin
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.GET("/", controllers.HandleWelcome)
	r.GET("/login", controllers.HandleLogin)
	r.GET("/callback", controllers.HandleCallback)
	r.GET("/search/:isrc", controllers.HandleSearchByISRC)
	// r.GET("/search/artist/:name", handleSearchByArtistName)

	fmt.Println("Server is listening on :8080...")
	r.Run(":8080")
}
