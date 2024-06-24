package main

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"net/http"
	"time"
	"github.com/DriveFluency/02-Backend/cmd/server/handler"
	"github.com/DriveFluency/02-Backend/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"github.com/DriveFluency/02-Backend/pkg/store"
	"github.com/DriveFluency/02-Backend/internal/pack"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/DriveFluency/02-Backend/docs"

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

	// mysql 
	db,err:= sql.Open("mysql","root:password@tcp(localhost:3306)/drive-fluency")
	if err != nil{
	  log.Fatal()
	  panic(err.Error())
  }

	// iniciar las entidades
	packStore := store.NewSqlPack(db)
	packRepo := pack.NewRepositoryPack(packStore)
	packService := pack.NewServicePack(packRepo)
	packHandler := handler.NewPackHandler(packService)

	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // TODO: Configuracion insegura, permite todos los origenes.
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "token"},
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

	roles := []string{"cliente", "admin"} 

	//endpoint packs 
	packs := r.Group("/packs")
	{
		packs.GET("/:id", packHandler.GetByID())
		packs.GET("", packHandler.GetAll())
		packs.POST("", middleware.AuthorizedJWT(roles), packHandler.Post())
		packs.PUT("/:id",middleware.AuthorizedJWT(roles), packHandler.Put())
		packs.PATCH("/:id",middleware.AuthorizedJWT(roles), packHandler.Patch())
		packs.DELETE("/:id",middleware.AuthorizedJWT(roles), packHandler.Delete())
		
	}

	r.POST("/login", handler.LoginHandler)
	r.POST("/logout", handler.LogoutHandler)
	r.POST("/register", handler.RegisterUserHandler)
	r.POST("/change", handler.ChangePasswordHandler)
	r.PUT("/profile", handler.UpdateProfile)
	r.GET("/reset", handler.ResetHandler)
	
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
	

	docs.SwaggerInfo.Host = "localhost:8085" //"http://conducirya.com.ar:8085" 
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8085")
}
