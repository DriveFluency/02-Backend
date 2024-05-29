package main

import (
	"02-Backend/pkg/middleware"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/coreos/go-oidc"
	"context"
	"log"
	"golang.org/x/oauth2"
)


var (
	clientID     = "drivefluency"
	clientSecret = "voJHDlSbYC69OzKDCA1CeGVkHXWhMxQd"
	redirectURI  = "http://localhost:8085/callback"
	realmURL  = "http://localhost:8090/realms/DriveFluency"
)

//callback 

var (
    oauth2Config = oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        Endpoint:     oauth2.Endpoint{
            AuthURL:  fmt.Sprintf("%s/protocol/openid-connect/auth", realmURL),
            TokenURL: fmt.Sprintf("%s/protocol/openid-connect/token", realmURL),
        },
        RedirectURL: redirectURI,
		Scopes:      []string{oidc.ScopeOpenID, "profile", "email"},
    }
)



func main() {


	r := gin.Default()

	r.GET("/login", loginHandler)
    r.GET("/callback", callbackHandler)

	roles:= []string{"usuario","admin"}
    r.Use(middleware.AuthorizedJWT(roles)) 

	
	endopointsPrueba := r.Group("/prueba")
	{
		endopointsPrueba.GET("/prueba",func (c *gin.Context){
				c.JSON(http.StatusOK, gin.H{"endpoint": "prueba"})
				return} )
}

r.Run(":8085")
}

func loginHandler(c *gin.Context) {
    redirectURL := fmt.Sprintf("%s/protocol/openid-connect/auth?client_id=%s&response_type=code&redirect_uri=%s", realmURL, clientID, redirectURI)
    c.Redirect(http.StatusFound, redirectURL)
}


// usuario autenticado obtiene el código de acceso

func callbackHandler(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not provided"})
        return
    }
	println(code)

    token, err := exchangeCodeForToken(code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code for token"+ err.Error()})
        return
    }

    
    c.SetCookie("access_token", token, 3600, "/", "localhost", false, true)


    c.Redirect(http.StatusFound, "http://localhost:8085/")
}


func exchangeCodeForToken(code string) (string, error) {
    ctx := context.Background()
    token, err := oauth2Config.Exchange(ctx, code)
    if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
        return "", err
    }

	log.Printf("Access token: %s", token.AccessToken)

    return token.AccessToken, nil
}