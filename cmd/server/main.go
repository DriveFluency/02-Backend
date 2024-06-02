package main

import (
	"github.com/DriveFluency/02-Backend/pkg/middleware"
    "github.com/DriveFluency/02-Backend/cmd/server/handler"
	"net/http"
    "github.com/gin-gonic/gin"
	"github.com/DriveFluency/02-Backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

	docs.SwaggerInfo.Host = "localhost:8085"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/login", handler.LoginHandler)
   // r.GET("/callback", handler.CallbackHandler)

	roles:= []string{"cliente","admin"}
    r.Use(middleware.AuthorizedJWT(roles)) 

	
	endopointsPrueba := r.Group("/prueba")
	{
		endopointsPrueba.GET("/",func (c *gin.Context){
			c.JSON(http.StatusOK, gin.H{"endpoint": "all users"})
				return} )

		endopointsPrueba.GET("/admin" ,func (c *gin.Context){
			c.JSON(http.StatusOK, gin.H{"endpoint": "only user admin"})
            return} )
}

r.Run(":8085")
}

