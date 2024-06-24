package util

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"errors"
	"strings"

	"github.com/DriveFluency/02-Backend/cmd/server/config"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"gopkg.in/resty.v1"
)

// ExtractTokenFromHeader extrae el token JWT del header Authorization del contexto gin.Context
func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	// Obtener el valor del header Authorization
	authHeader := c.GetHeader("Authorization")

	// Verificar si el header Authorization está presente
	if authHeader == "" {
		return "", errors.New("missing auth token")
	}

	// Verificar si comienza con "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid token format")
	}

	// Extraer el token eliminando "Bearer " del inicio
	token := authHeader[len("Bearer "):]

	return token, nil
}

func ExtractUserIDFromJWT(rawAccessToken string, cfg *config.Config) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Duration(10) * time.Second,
		Transport: tr,
	}

	ctx := oidc.ClientContext(context.Background(), client)

	provider, err := oidc.NewProvider(ctx, cfg.RealmURL)
	if err != nil {
		return "", fmt.Errorf("error al obtener el proveedor OIDC: %v", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	idToken, err := verifier.Verify(ctx, rawAccessToken)
	if err != nil {
		return "", fmt.Errorf("error al verificar el token JWT: %v", err)
	}

	claims := map[string]interface{}{}
	if err := idToken.Claims(&claims); err != nil {
		return "", fmt.Errorf("error al extraer las reclamaciones del token: %v", err)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("no se encontró el campo 'sub' en las reclamaciones del token")
	}

	return userID, nil
}

func GetAdminUserToken(cfg *config.Config) (string, error) {
	// TODO: Es inseguro usar Direct Access Grants.
	// Se recomienda cambiarlo para user ClientID y Client Secret
	resp, err := resty.R().SetFormData(map[string]string{
		"client_id":  cfg.AdminClientID,
		"username":   cfg.AdminUser,
		"password":   cfg.AdminPass,
		"grant_type": "password",
		"scope":      "openid",
	}).Post(cfg.TokenURL)

	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", errors.New("access_token not found in response")
	}

	return token, nil
}
