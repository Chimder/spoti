package httpgin

import (
	// _ "csTrade/docs"
	// "csTrade/internal/handlers/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// swaggerfiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// offerServ := service.NewOfferService(repo, botmanager)
	// offerHandler := NewOfferHandler(offerServ)

	// userServ := service.NewUserService(repo)
	// userHandler := NewUserHandler(userServ)

	// transactionServ := service.NewTransactionService(repo)
	// transactionHandler := NewTransactionHandler(transactionServ)

	{
		// r.GET("/swagger", ginSwagger.WrapHandler(swaggerfiles.Handler))
		// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		r.GET("/healthz", func(c *gin.Context) {
			c.String(200, "ok")
		})
	}

	// api := r.Group("/api/v1")

	// api.POST("/users/create", userHandler.CreateUser)

	// users := api.Group("/users").Use(middleware.AuthMiddleware())
	// {
	// 	users.GET("/:id")
	// 	users.GET("/:id/cash")
	// 	users.PATCH("/:id/cash")
	// }


	return r
}
