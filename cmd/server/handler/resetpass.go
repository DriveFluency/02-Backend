package handler

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)



func ResetHandler(c *gin.Context) {
resetURL := fmt.Sprintf("%s/login-actions/reset-credentials?client_id=%s&redirect_uri=http://localhost:8085/login",realmURL,clientID) //&redirect_uri=http://conducirya.com.ar
c.Redirect(http.StatusFound, resetURL)
}