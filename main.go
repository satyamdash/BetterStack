package main

import (
	"fmt"
	"time"

	"github.com/BetterStack/producer"
	redisclient "github.com/BetterStack/redis"
	"github.com/BetterStack/routers"
	"github.com/BetterStack/worker"
	"github.com/gin-gonic/gin"
)

func main() {
	redisclient.Init("localhost:6379", "", 0)
	router := gin.Default()

	router.POST("/user/signin", routers.SigninRouter)
	router.POST("/user/signup", routers.SignupRouter)

	protected := router.Group("/api/v1")

	protected.Use(AuthMiddleware())
	{
		protected.GET("/website", routers.WebsiteAdd)
		protected.GET("/status/:websiteId", routers.GetStatus)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Worker crashed: %v\n", r)
			}
		}()
		worker.Worker()
	}()

	go func() {
		// Run producer immediately
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Producer crashed: %v\n", r)
				}
			}()
			producer.Producer()
		}()

		// Then run on ticker
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Producer crashed: %v\n", r)
					}
				}()
				producer.Producer()
			}()
		}
	}()

	router.Run(":3000")

}
