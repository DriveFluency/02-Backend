package handler

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/DriveFluency/02-Backend/cmd/server/config"
	"github.com/DriveFluency/02-Backend/cmd/server/util"

	"github.com/gin-gonic/gin"
	"gopkg.in/resty.v1"
)

type KeycloakConfig struct {
	BaseURL      string
	AdminUser    string
	AdminPass    string
	ClientID     string
	ClientSecret string
	Realm        string
}

type Profile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	DNI       string `json:"dni"`
	Telefono  string `json:"telefono"`
	Email     string `json:"email"`
	Ciudad    string `json:"ciudad"`
	Localidad string `json:"localidad"`
	Direccion string `json:"direccion"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// Mapa para traducir nombres de campos de inglés a español
var spanishFieldNames = map[string]string{
	"FirstName": "Nombre",
	"LastName":  "Apellido",
	"DNI":       "DNI",
	"Telefono":  "Teléfono",
	"Email":     "Correo electrónico",
	"Ciudad":    "Ciudad",
	"Localidad": "Localidad",
	"Direccion": "Dirección",
	"Username":  "Nombre de usuario",
	"Password":  "Contraseña",
}

func RegisterUserHandler(c *gin.Context) {
	cfg := config.GetConfig()
	var profile Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	missingFields, err := validateProfile(profile)
	if err != nil {
		// Construir mensaje de error con nombres de campos en español
		var spanishMissingFields []string
		for _, fieldName := range missingFields {
			if spanishName, ok := spanishFieldNames[fieldName]; ok {
				spanishMissingFields = append(spanishMissingFields, spanishName)
			} else {
				spanishMissingFields = append(spanishMissingFields, fieldName)
			}
		}

		errorMessage := "Los siguientes campos son requeridos pero están vacíos: " + strings.Join(spanishMissingFields, ", ")
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	token, err := util.GetAdminUserToken(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en registro (imposible obtener token admin): " + err.Error()})
		return
	}

	user := map[string]interface{}{
		"username":  profile.Username,
		"email":     profile.Email,
		"firstName": profile.FirstName,
		"lastName":  profile.LastName,
		"enabled":   true,
		"attributes": map[string]string{
			"DNI":       profile.DNI,
			"telefono":  profile.Telefono,
			"ciudad":    profile.Ciudad,
			"localidad": profile.Localidad,
			"direccion": profile.Direccion,
		},
		"credentials": []map[string]string{
			{
				"type":      "password",
				"value":     profile.Password,
				"temporary": "false",
			},
		},
	}

	resp, err := resty.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(user).
		Post(cfg.UserURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en registro (error en keycloack): " + err.Error()})
		return
	}

	if resp.StatusCode() >= 400 {
		c.JSON(resp.StatusCode(), gin.H{"error": "Error en registro: " + string(resp.Body())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registro exitoso"})
}

func validateProfile(profile Profile) ([]string, error) {
	// Utilizar el paquete reflect para obtener los tags de estructura y validar
	v := reflect.ValueOf(profile)
	vt := v.Type()

	var missingFields []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := vt.Field(i).Tag.Get("json")

		if field.Interface() == "" {
			missingFields = append(missingFields, tag)
		}
	}

	return missingFields, nil
}
