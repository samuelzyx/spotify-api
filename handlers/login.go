package handlers

import (
	"net/http"

	"spotify-api/config"

	"github.com/gin-gonic/gin"
)

var state = "state-string"

func HandleLogin(c *gin.Context) {
	url := config.OAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}
