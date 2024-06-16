package handler

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"bytes"
	"encoding/json"
	"log"
	"io"
	"crypto/tls"
	"time"
)


type Request struct{
	Username string `json:"username" binding:"required"`
}

var UrlResetPrueba =fmt.Sprintf("%s/login-actions/reset-credentials?client_id=%s",realmURL,clientID)
//con auth despues 

func ResetHandlerRedirect(c *gin.Context){
	  //&redirect_uri=http://conducirya.com.ar//&redirect_uri=http://localhost:8085/login
	c.Redirect(http.StatusFound, UrlResetPrueba)
}

func ResetHandler(c *gin.Context) {
	
resetURL := UrlResetPrueba
var req Request
err := c.ShouldBindJSON(&req) // lo pasa a estructura deja de ser json 
if err != nil{
	log.Println("Invalid request:", err)
	c.JSON(http.StatusBadRequest,gin.H{"error":"invalid request"})
	return 
}

// carga útil a enviar clave-valor guarda el json en el formato esperado por keycloak 
body:= map[string]string{
	"username": req.Username, 
}


// convierte a JSON 
bodyBytes, err := json.Marshal(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}

//creando la solicitud  
request, err := http.NewRequest("POST",resetURL, bytes.NewBuffer(bodyBytes)) //convierte una cadena de bytes a obj de tipo io.reader
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
	return
}
// aclaro que envio un JSON
request.Header.Set("Content-Type","application/json")

/* Obtener las cookies de la solicitud original*/
for _, cookie := range c.Request.Cookies() {
	request.AddCookie(cookie)
}


log.Println("print request:", request)


tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{
			Timeout:   10 * time.Second,
			Transport: tr,
		}

response,err:=client.Do(request)
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
	return
}

defer response.Body.Close()

log.Println("response:", response)

// vemos cual es error puntual 
responseKeycloak,err:= io.ReadAll(response.Body)
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
	return
}
log.Println("Lo que responde keycloak e interpreta go:", string(responseKeycloak))

//204 no content 
if response.StatusCode == http.StatusNoContent || response.StatusCode == http.StatusOK{
	c.JSON(http.StatusOK, gin.H{"message":"Mensaje de recupero enviado a su correo"})
	return
}

//lo guardamos en una interfaz clave valor... 
var mapResponseKeycloak map[string]interface{}

err = json.Unmarshal(responseKeycloak, &mapResponseKeycloak) // transforma el JSON recibido pero el problema es q recibe html error cookie no encontrada
	if err != nil {
		log.Println("error unmarshall:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}
log.Printf("error en la respuesta de keycloak, %v",mapResponseKeycloak )
c.JSON(response.StatusCode, mapResponseKeycloak)

}


	
