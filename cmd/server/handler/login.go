package handler


import (
	"fmt"
	"net/http"
    "github.com/gin-gonic/gin"
	//"github.com/coreos/go-oidc"
	"context"
	"log"
	"golang.org/x/oauth2"
)

var (
	clientID     = "drivefluency"
	clientSecret = "voJHDlSbYC69OzKDCA1CeGVkHXWhMxQd"
	realmURL  = "http://localhost:8090/realms/DriveFluency" // cambiar 
   // redirectURI  = "http://localhost:8085/callback"
    tokenURL     = fmt.Sprintf("%s/protocol/openid-connect/token", realmURL)
	authURL      = fmt.Sprintf("%s/protocol/openid-connect/auth", realmURL)
    

)

type RequestBody struct{

    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`

}



// @Summary Inicio de sesión
// @Tag Login
//@Accept json
// @Produce json
// @Success 200
// @Router /login [post]
func LoginHandler(c *gin.Context){

    var body RequestBody
    err := c.ShouldBindJSON(&body)
    if err != nil {
       c.JSON(http.StatusBadRequest,gin.H{"error":"invalid request"} )
       return 
     }

     // obtener token 
     token, err := authenticateUser(body.Username, body.Password)
     if err != nil {
         c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
         return
     }
 
     c.SetCookie("access_token", token, 3600, "/", "localhost", false, true)

     c.JSON(http.StatusOK, gin.H{"access_token": token})

     // redirigir al home del front ? 
    // c.Redirect(http.StatusFound, url del home del front )


}

//postman
func authenticateUser(username, password string) (string, error) {
	oauth2Config := &oauth2.Config{
 	ClientID:     clientID,
 	ClientSecret: clientSecret,
 	Endpoint: oauth2.Endpoint{
 		TokenURL: tokenURL,
 		AuthURL:  authURL,
 	},
 	Scopes: []string{"roles","email"}, //"profile", "email",
 }

	ctx := context.Background()
	token, err := oauth2Config.PasswordCredentialsToken(ctx,username,password)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		return "", err
	}

	return token.AccessToken, nil
}




/*//callback 

var (
    oauth2Config = oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        Endpoint:     oauth2.Endpoint{
            AuthURL:  fmt.Sprintf("%s/protocol/openid-connect/auth", realmURL),
            TokenURL: fmt.Sprintf("%s/protocol/openid-connect/token", realmURL),
        },
        RedirectURL: redirectURI,
		Scopes:      []string{oidc.ScopeOpenID,"roles","email" }, // revisar "profile", "email"

    }
)


func LoginHandler(c *gin.Context) {
    redirectURL := fmt.Sprintf("%s/protocol/openid-connect/auth?client_id=%s&response_type=code&redirect_uri=%s&scope=roles email", realmURL, clientID, redirectURI)
    c.Redirect(http.StatusFound, redirectURL)
}


// usuario autenticado obtiene el código de acceso

func CallbackHandler(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not provided"})
        return
    }
	println("este es el codigo de autorización "+code)

    token, err := exchangeCodeForToken(code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code for token"+ err.Error()})
        return
    }

    
    c.SetCookie("access_token", token, 3600, "/", "localhost", false, true)




    c.Redirect(http.StatusFound, "http://localhost:8085/prueba")
}


func exchangeCodeForToken(code string) (string, error) { 
    ctx := context.Background() 
    token, err := oauth2Config.Exchange(ctx, code) // revisar
    if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
        return "", err
    }

	log.Printf("Access token: %s", token.AccessToken)

    return token.AccessToken, nil
}

*/