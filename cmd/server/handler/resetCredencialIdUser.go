package handler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"time"
)

var UserURL="http://conducirya.com.ar:18080/admin/realms/DriveFluency/users" 

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// genera el token para el cliente 
func getAdminToken() (string, error) {
	
	data := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", clientID, clientSecret)
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{
		Timeout:   time.Duration(10000) * time.Second,
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	log.Println("token de respuesta al identificador del cliente", tokenResponse)
	return tokenResponse.AccessToken, nil
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}


//figura que  no tiene permisos para listar los usuarios 
func findUserID(username, token string) (string, error) {
	userURL := fmt.Sprintf("%s?username=%s",UserURL, username)

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{
		Timeout:   time.Duration(10000) * time.Second,
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	log.Println("usuarios con ese nombre retornados antes de pasar a bytes", resp)
	 // aca arroja un 403 , no se tienen los permisos para acceder a los usuarios
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		log.Println("error al guardar en un slice los usuarios ", err)
		return "", err
	}

	log.Println("usuarios ya guardados en slice", users)

	if len(users) == 0 {
		return "", fmt.Errorf("user not found")
	}

	return users[0].ID, nil
}

func ResetHandler2(c *gin.Context) {
	var req R
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := getAdminToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get admin token"})
		return
	}

	userID, err := findUserID(req.Username, token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	resetPasswordURL := fmt.Sprintf("%s/%s/execute-actions-email", UserURL, userID)
	reqBody := []string{"UPDATE_PASSWORD"}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}

	request, err := http.NewRequest("PUT", resetPasswordURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{
		Timeout:   time.Duration(10000) * time.Second,
		Transport: tr,
	}

	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent || response.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"message": "Reset password email sent"})
	} else {
		body, _ := io.ReadAll(response.Body)
		c.JSON(response.StatusCode, gin.H{"error": string(body)})
	}
}
