package main

import (
	"net/http"
	"time"

	"github.com/DriveFluency/02-Backend/cmd/server/handler"
	"github.com/DriveFluency/02-Backend/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// TODO: Habilitar cuando esten los endpoints listos
	// "github.com/DriveFluency/02-Backend/docs"
	// "github.com/swaggo/files"
	// "github.com/swaggo/gin-swagger"
)

// @title Drive Fluency
// @version 1.0
// @description This API Handle Shifts.
// @termsOfService https://terminos-y-condiciones

// @contact.name API Support
// @contact.url

// @license.name
// @license.url

func main() {
	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // TODO: Configuracion insegura, permite todos los origenes.
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	/*
		// TODO: Habilitar cuando esten los endpoints listos
		docs.SwaggerInfo.Host = "localhost:8085"
		docs.SwaggerInfo.Title = "Drive Fluency API"
		docs.SwaggerInfo.Description = "Drive Fluency API"
		docs.SwaggerInfo.Version = "1.0"
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	*/

	r.GET("/reset", handler.ResetHandler)

	r.POST("/login", handler.LoginHandler)
	r.POST("/logout", handler.LogoutHandler)
	r.POST("/register", handler.RegisterUserHandler)

	roles := []string{"cliente", "admin"}

	endopointsPrueba := r.Group("/prueba")
	endopointsPrueba.Use(middleware.AuthorizedJWT(roles))
	{
		endopointsPrueba.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"endpoint": "all users"})
		})

		endopointsPrueba.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"endpoint": "only user admin"})
		})
	}

	r.Run(":8085")
}
