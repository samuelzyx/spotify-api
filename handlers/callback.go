package handlers

import (
	"context"
	"net/http"
	"spotify-api/config"

	"github.com/gin-gonic/gin"
)

var authenticatedClient *http.Client

func HandleCallback(c *gin.Context) {
	code := c.Query("code")
	tok, err := config.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error exchanging code for access token")
		return
	}

	authenticatedClient = config.OAuthConfig.Client(context.Background(), tok)
	c.String(http.StatusOK, "Successfully authenticated with Spotify. Access token: %s", tok.AccessToken)
}
