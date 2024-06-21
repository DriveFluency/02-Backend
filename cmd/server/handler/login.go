package handler

import (
	"fmt"
	"net/http"
	"github.com/DriveFluency/02-Backend/internal"
	"github.com/gin-gonic/gin"
	"github.com/coreos/go-oidc"
	"context"
	"log"
	"golang.org/x/oauth2"
	"time"
	"crypto/tls"
)

var (
	clientID     = "drivefluency"
	clientSecret = "083E22w85Iw9T2vctotLkT3ZAEDaqXsA"                   //
	realmURL     = "http://conducirya.com.ar:18080/realms/DriveFluency" //http://localhost:8090/realms/DriveFluency"
	tokenURL = fmt.Sprintf("%s/protocol/openid-connect/token", realmURL)
	authURL  = fmt.Sprintf("%s/protocol/openid-connect/auth", realmURL)
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

	/**/datosUsuario, err := saveUser(token, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "save properties of user"})
		return
	}
	//c.SetCookie("access_token", token, 3600, "/", "http://conducirya.com.ar", false, true) // ver el dominio en el cual estaría habilitada
	c.Header("access_token", token)
	c.JSON(http.StatusOK, gin.H{"access_token": token, "user": datosUsuario}) //

	// Ya lo redireccionan del front
	//c.Redirect(http.StatusFound, "http://conducirya.com.ar")

}

// postman
func authenticateUser(username, password string) (string, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
			AuthURL:  authURL,
		},
		Scopes: []string{"roles","usuario"}, //"profile", "email",DNI con atribbute
	}

	ctx := context.Background()
	token, err := oauth2Config.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		return "", err
	}
	return token.AccessToken, nil
}

func saveUser(token string , c *gin.Context) (internal.User,error){

  var Usuario internal.User

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   time.Duration(10000) * time.Second,
		Transport: tr,
	}

	ctx := oidc.ClientContext(context.Background(), client)

	provider, err := oidc.NewProvider(ctx, realmURL)
	if err != nil {
		log.Printf("Error getting the provider: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed save user while getting the provider"})	
		
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	verifier := provider.Verifier(oidcConfig)
	log.Printf("devuelve un IDTokenVerifier que utiliza el conjunto de claves del proveedor para verificar los JWT.")

	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify token"})
	
	}

	if err := idToken.Claims(&Usuario); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to extract claims"})
	}
	log.Print("esto es lo que trajo del usuario",Usuario)

	return Usuario, nil


}
