package handler

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

// Utiliza las mismas variables del paquete que estan en login.go
var (
	logoutURL = fmt.Sprintf("%s/protocol/openid-connect/logout", realmURL)
)

func LogoutHandler(c *gin.Context) {

	//buscamos el access token en la cookie
	 //rawAccessToken,err := c.Cookie("access_token")    
	 rawAccessToken:= c.GetHeader("token")
	log.Printf("Raw Access token obteniendolo del header endpoint log out: %s", rawAccessToken)
	if rawAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access token not found"})
		return
	}

	// crea un cliente y no valida los certificados TLS https
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   time.Duration(10000) * time.Second,
		Transport: tr,
	}

	// inicia el contexto oidc para realizar solicitudes desde el cliente http
	ctx := oidc.ClientContext(context.Background(), client)

	provider, err := oidc.NewProvider(ctx, realmURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed log out while getting the provider"})
		return
	}

	// espero que traiga el token con las claves de validacion para el cliente drivefluency
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)
	log.Printf("devuelve un IDTokenVerifier que utiliza el conjunto de claves del proveedor para verificar los JWT.")

	idToken, err := verifier.Verify(ctx, rawAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify token"})
		return
	}

	log.Printf("token de validación verificado por el IAM: %s", idToken)

	//func (*oidc.IDToken).Claims(v interface{}) error Claims unmarshals the raw JSON payload of the ID Token into a provided struct
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to extract claims"})
		return
	}

	// JTI del token
	jti, ok := claims["jti"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JTI not found in token claims"})
		return
	}

	log.Printf("este es el jti que trae %s", jti)

	// prepara la solicitud de cierre de sesión a keycloak asociandolo al contexto definido
	request, err := http.NewRequestWithContext(ctx, "POST", logoutURL, strings.NewReader(fmt.Sprintf("client_id=%s&jti=%s", clientID, jti)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create logout request"})
		return
	}

	// con esto indicamos a keycloak que el cuerpo de la solicitud esta codificado como un formulario URL
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// envio de la solicitud con metodo do de http
	logOut, err := client.Do(request)
	if err != nil || logOut.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed logout"})
		return
	}
	log.Printf("cierre de sesión exitoso %s", logOut.Status)

	//borrar la cookie y redireccionar al home del front
	//c.SetCookie("access_token", "", -1, "/", "http://conducirya.com.ar", false, true)
	
	// el local storage se borra del lado del front... 

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	//c.Redirect(http.StatusFound, "http://conducirya.com.ar")

}
