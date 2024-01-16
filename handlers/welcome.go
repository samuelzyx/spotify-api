package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWelcome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
