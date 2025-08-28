package routers

import (
	"os"
	"time"

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
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(), // 1 hour expiry
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"status": "signin successful", "token": tokenString})
}
