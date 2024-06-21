package middleware

import (
	"context"
	"crypto/tls"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	// "fmt"
)

var (
	RealmConfigURL string =  "http://conducirya.com.ar:18080/realms/DriveFluency"//si agregas :18080 te dice q estpera esta pero asi es un 404 "http://localhost:8090/realms/DriveFluency"
	clientID       string = "drivefluency"
	// redirectURI  = "http://localhost:8085/callback" // no va )
)

type Res401Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"HttpCode" example:"401"`
	Message  string `json:"message" example:"authorization failed"`
}

type Claims struct {
	ResourceAccess client `json:"resource_access,omitempty"`
	JTI            string `json:"jti,omitempty"`
}

type client struct {
	DriveFluency clientRoles `json:"DriveFluency,omitempty"`
}

type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

func AuthorizedJWT(roles []string) gin.HandlerFunc {

	return func(c *gin.Context) {

		//deberian funcionar ambas
		//	rawAccessToken := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
		rawAccessToken := c.GetHeader("token") //c.Cookie("access_token")
		log.Printf("Raw Access token en el middleware: %s", rawAccessToken)

		/* cambiar redirigir al endpoint login */
		if rawAccessToken == "" {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "access token not found"})
			c.Redirect(http.StatusFound, "http://conducirya.com.ar/SignIn") //ver si esta bien la ruta	
			//c.Abort()
			return
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{
			Timeout:   time.Duration(10000) * time.Second,
			Transport: tr,
		}

		ctx := oidc.ClientContext(context.Background(), client)

		provider, err := oidc.NewProvider(ctx, RealmConfigURL)
		if err != nil {
			authorizationFailed("authorization failed while getting the provider: "+err.Error(), c)
			return
		}

		// espero que traiga el token con las claves de validacion para el cliente drivefluency pero retorna para el cliente account
		oidcConfig := &oidc.Config{
			ClientID: clientID,
		}
		verifier := provider.Verifier(oidcConfig)
		log.Printf("devuelve un IDTokenVerifier que utiliza el conjunto de claves del proveedor para verificar los JWT.")

		idToken, err := verifier.Verify(ctx, rawAccessToken)
		if err != nil {
			authorizationFailed("authorization failed while verifying the token: "+err.Error(), c)
			return
		}

		var IDTokenClaims Claims
		if err := idToken.Claims(&IDTokenClaims); err != nil {
			authorizationFailed("claims : "+err.Error(), c)
			return
		}

		user_access_roles := IDTokenClaims.ResourceAccess.DriveFluency.Roles
		for _, b := range user_access_roles {
			log.Printf("ROL %s", b)
			if b == roles[0] && c.FullPath() == "/prueba/" {
				c.Next()
				return
			}
			if b == roles[1] {
				c.Next()
				return
			}
		}
		authorizationFailed("user not allowed to access this api", c)
	}
}

func authorizationFailed(message string, c *gin.Context) {
	data := Res401Struct{
		Status:   "FAILED",
		HTTPCode: http.StatusUnauthorized,
		Message:  message,
	}

	c.AbortWithStatusJSON(200, gin.H{"response": data})

}
