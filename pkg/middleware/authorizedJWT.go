package middleware

import(

	 "strings"
	 "github.com/gin-gonic/gin"
	 "net/http" 
	 "time"
	 "github.com/coreos/go-oidc"
	 "context"
	 "crypto/tls"
	 "fmt"
	)


var (
	RealmConfigURL string = "http://localhost:8090/realms/DriveFluency"
    clientID string = "drivefluency"
    redirectURI  = "http://localhost:8085/callback")


type Res401Struct struct{

	Status string `json:"status" example:"FAILED"`
	HTTPCode int `json:"HttpCode" example:"401"`
	Message string `json:"message" example:"authorization failed"`
}


type Claims struct {
	ResourceAccess client `json:"resource_access,omitempty"`
	JTI string `json:"jti,omitempty"`
}

type client struct {
	DriveFluency clientRoles `json:"DriveFluency,omitempty"`
}


type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

func AuthorizedJWT(roles []string) gin.HandlerFunc {

	return func(c *gin.Context) {

		
		rawAccessToken := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1) 

		if rawAccessToken == "" {
			redirectURL := fmt.Sprintf("%s/protocol/openid-connect/auth?client_id=%s&response_type=code&redirect_uri=%s", RealmConfigURL, clientID, redirectURI)
    		c.Redirect(http.StatusFound, redirectURL)
			return 
		}


	
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}

		ctx := oidc.ClientContext(context.Background(), client)

		provider, err := oidc.NewProvider(ctx, RealmConfigURL)
		if err != nil {
			authorizationFailed("authorization failed while getting the provider: "+err.Error(), c)
			return
		}
		oidcConfig := &oidc.Config{
			ClientID: clientID,
		}

		
		verifier := provider.Verifier(oidcConfig)
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
			
			if b == roles[0] && c.FullPath() == "/prueba" { 
				c.Next() 
				return
			}
			if b == roles[1]{
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
