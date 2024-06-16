package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/clientcredentials"
)

type R struct {
    Username string `json:"username" binding:"required"`
}

func getClientCredentialsToken(ctx context.Context, clientID, clientSecret, realm string) (string, error) {
   
	providerURL := fmt.Sprintf("http://localhost:8090/realms/%s", realm)

    provider, err := oidc.NewProvider(ctx, providerURL)
    if err != nil {
		log.Println("error al obtener el proveedor ",err)
        return "", fmt.Errorf("failed to get provider: %v", err)
    }

    oauth2Config := clientcredentials.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        TokenURL:     provider.Endpoint().TokenURL,
    }

    token, err := oauth2Config.Token(ctx)
    if err != nil {
		log.Println("error al obtener el token del cliente ",err)
        return "", fmt.Errorf("failed to get token: %v", err)
    }
    return token.AccessToken, nil
}

func ResetPasswordHandler(c *gin.Context) {
    var req Request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    clientID := "drivefluency"
    clientSecret := "UMQuQX26AD63348ftkzL8c2AyBB05s3f"
    realm := "DriveFluency"
    //ctx := context.Background()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{
		Timeout:   time.Duration(10000) * time.Second,
		Transport: tr,
	}

	ctx := oidc.ClientContext(context.Background(), client)

    token, err := getClientCredentialsToken(ctx, clientID, clientSecret, realm)
    if err != nil {
		log.Println("error al obtener el token del cliente ",err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
        return
    }

    resetPasswordURL := fmt.Sprintf("http://localhost:8090/realms/%s/login-actions/reset-credentials?client_id=%s", realm, clientID)
    body := map[string]string{
        "username": req.Username,
    }
    bodyBytes, err := json.Marshal(body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
        return
    }

    request, err := http.NewRequest("POST", resetPasswordURL, bytes.NewBuffer(bodyBytes))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
        return
    }

    request.Header.Set("Authorization", "Bearer "+token)
    request.Header.Set("Content-Type", "application/json")

	/*//agrego cookies
	for _, cookie := range c.Request.Cookies() {
		request.AddCookie(cookie)}
	*/

    client2 := &http.Client{}
    response, err := client2.Do(request)
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


