package handler

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
)

type changePass struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

func ChangePasswordHandler(c *gin.Context) {
	var req changePass
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	rawAccessToken := c.GetHeader("token")
	log.Printf("Raw Access token obteniendolo del header endpoint log out: %s", rawAccessToken)
	if rawAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access token not found"})
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Duration(10000) * time.Second,
		Transport: tr,
	}

	ctx := oidc.ClientContext(context.Background(), client)

	provider, err := oidc.NewProvider(ctx, realmURL)
	if err != nil {
		log.Printf("Error getting the provider: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get provider"})
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)
	log.Printf("devuelve un IDTokenVerifier que utiliza el conjunto de claves del proveedor para verificar los JWT.")

	idToken, err := verifier.Verify(ctx, rawAccessToken)
	if err != nil {
		log.Printf("Error verifying token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify token"})
		return
	}
	log.Println(idToken)

	claims := map[string]interface{}{}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to extract claims"})
		return
	}

	// guardo el id del usuario obtenido desde el access token 
	userID := claims["sub"].(string)
	adminToken, err := GetAdminToken()
	if err != nil {
		log.Printf("Error getting admin token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get admin token"})
		return
	}

	err = changeUserPassword(adminToken, userID, req.NewPassword)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
func changeUserPassword(adminToken, userID, newPassword string) error {
	changePasswordURL := fmt.Sprintf("http://conducirya.com.ar:18080/admin/realms/DriveFluency/users/%s/reset-password", userID)
	client := &http.Client{}
	reqBody := map[string]interface{}{
		"type":"password",
		"temporary":false,
		"value":newPassword,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("PUT", changePasswordURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("failed to change password, status code: %d, response: %s", resp.StatusCode, bodyString)
	}

	return nil
}
