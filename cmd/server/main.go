package main

import (
	"net/http"
	"time"
	"github.com/DriveFluency/02-Backend/cmd/server/handler"
	//"github.com/DriveFluency/02-Backend/docs"
	"github.com/DriveFluency/02-Backend/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
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
		AllowOrigins:     []string{"*"},//"http://conducirya.com.ar"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length","Authorization"},
		AllowCredentials: true,
		/*AllowOriginFunc: func(origin string) bool {
			return origin == "http://conducirya.com.ar"
		},*/
		MaxAge: 12 * time.Hour,
	}))

	/*docs.SwaggerInfo.Host = "localhost:8085"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))*/

	r.POST("/login", handler.LoginHandler)
	// r.GET("/callback", handler.CallbackHandler)
	r.POST("/logout", handler.LogoutHandler)
	r.GET("/reset", handler.ResetHandler)

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
