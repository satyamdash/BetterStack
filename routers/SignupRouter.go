package routers

import (
	"github.com/BetterStack/db"
	"github.com/gin-gonic/gin"
)

type SignUp struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignupRouter(c *gin.Context) {
	var client = db.GetPrismaClient()
	var json SignUp

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(403, gin.H{"error": ""})
		return
	}
	// Check if user already exists

	usr, err := client.User.CreateOne(
		db.User.Username.Set(json.User),
		db.User.Password.Set(json.Password),
	).Exec(c)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(200, gin.H{"status": "signup successful", "user": usr.Username})
}
