package handler

import (
	"fmt"
	"net/http"

	"github.com/DriveFluency/02-Backend/cmd/server/config"
	"github.com/DriveFluency/02-Backend/cmd/server/util"
	"github.com/gin-gonic/gin"
	"gopkg.in/resty.v1"
)

type Attributes struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	DNI       string `json:"dni"`
	Telefono  string `json:"telefono"`
	Ciudad    string `json:"ciudad"`
	Localidad string `json:"localidad"`
	Direccion string `json:"direccion"`
}

func UpdateProfile(c *gin.Context) {
	cfg := config.GetConfig()

	var profile Attributes

	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := util.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	adminToken, err := util.GetAdminUserToken(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en registro (imposible obtener token admin): " + err.Error()})
		return
	}

	userID, err := util.ExtractUserIDFromJWT(accessToken, cfg)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userProfile := map[string]interface{}{
		"username":  profile.Email,
		"email":     profile.Email,
		"firstName": profile.FirstName,
		"lastName":  profile.LastName,
		"attributes": map[string]string{
			"DNI":       profile.DNI,
			"telefono":  profile.Telefono,
			"ciudad":    profile.Ciudad,
			"localidad": profile.Localidad,
			"direccion": profile.Direccion,
		},
	}

	resp, err := resty.R().
		SetAuthToken(adminToken).
		SetHeader("Content-Type", "application/json").
		SetBody(userProfile).
		Put(fmt.Sprintf("%s/%s", cfg.UserURL, userID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la solicitud: " + err.Error()})
		return
	}

	if resp.StatusCode() >= 400 {
		c.JSON(resp.StatusCode(), gin.H{"error": "Error en la solicitud: " + string(resp.Body())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registro exitoso"})
}
