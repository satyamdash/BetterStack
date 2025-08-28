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
		protected.GET("/website", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Welcome to your website"})
		})
		protected.GET("/status/:websiteId", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "This is the status page for website " + c.Param("websiteId")})
		})
	}

	router.Run(":3000")

}
