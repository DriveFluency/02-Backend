package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	clientID     = "drivefluency"
	clientSecret = "083E22w85Iw9T2vctotLkT3ZAEDaqXsA"
	realmURL     = "http://conducirya.com.ar:18080/realms/DriveFluency"
	tokenURL     = fmt.Sprintf("%s/protocol/openid-connect/token", realmURL)
	authURL      = fmt.Sprintf("%s/protocol/openid-connect/auth", realmURL)
)

type RequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Inicio de sesión
// @Tag Login
// @Accept json
// @Produce json
// @Success 200
// @Router /login [post]
func LoginHandler(c *gin.Context) {

	var body RequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	// obtener token, este es el error q tiene el front
	token, err := authenticateUser(body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	log.Printf("Access Token:  %s\n", token)
	c.Header("access_token", token)
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func authenticateUser(username, password string) (string, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
			AuthURL:  authURL,
		},
	}

	ctx := context.Background()
	token, err := oauth2Config.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		return "", err
	}

	//
	return token.AccessToken, nil
}
