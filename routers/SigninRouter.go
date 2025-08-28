package routers

import (
	"github.com/BetterStack/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Signin struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SigninRouter(c *gin.Context) {
	client := db.GetPrismaClient()
	var json Signin

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(403, gin.H{"error": ""})
		return
	}

	user, err := client.User.FindFirst(
		db.User.Username.Equals(json.User),
	).Exec(c)

	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	if user.Password != json.Password {
		c.JSON(403, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Username,
	})

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"status": "signin successful", "token": tokenString})
}
