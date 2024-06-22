package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResetHandler(c *gin.Context) {
	resetURL := fmt.Sprintf("%s/login-actions/reset-credentials?client_id=%s&redirect=http://conducirya.com.ar", realmURL, clientID)
	c.Redirect(http.StatusFound, resetURL)
}
