package routers

import (
	"fmt"
	"time"

	"github.com/BetterStack/db"
	"github.com/gin-gonic/gin"
)

type Website struct {
	URL string `json:"url" binding:"required"`
}

func WebsiteAdd(c *gin.Context) {
	var website Website
	if err := c.ShouldBindJSON(&website); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	client := db.GetPrismaClient()
	userId := c.MustGet("userId").(string)
	fmt.Println("---------------------------------------")
	fmt.Println(userId)
	url, err := client.Website.CreateOne(
		db.Website.URL.Set(website.URL),
		db.Website.TimeAdded.Set(time.Now()),
		db.Website.User.Link(
			db.User.ID.Equals(userId),
		),
	).Exec(c)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Welcome to your website", "url": url.URL, "userId": userId})
}
