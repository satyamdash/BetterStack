package main

import (
	"github.com/BetterStack/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/user/signin", routers.SigninRouter)
	router.POST("/user/signup", routers.SignupRouter)

	router.Run(":3000")

}
