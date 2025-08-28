package main

import (
	"github.com/BetterStack/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/user/signin", routers.SigninRouter)
	router.POST("/user/signup", routers.SignupRouter)

	protected := router.Group("/api/v1")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/website", routers.WebsiteAdd)
		protected.GET("/status/:websiteId", routers.GetStatus)
	}

	router.Run(":3000")

}
