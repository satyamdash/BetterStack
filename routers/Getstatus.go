package routers

import (
	"github.com/BetterStack/db"
	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	websiteId := c.Param("websiteId")
	userId := c.MustGet("userId").(string)

	client := db.GetPrismaClient()

	website, err := client.Website.FindFirst(
		db.Website.ID.Equals(websiteId),
		db.Website.UserID.Equals(userId), // filter by logged-in user
	).With(
		db.Website.Ticks.Fetch().OrderBy(
			db.WebsiteTick.CreatedAt.Order(db.DESC), // order ticks by createdAt desc
		).Take(1), // only latest tick
	).Exec(c)

	if err != nil {
		c.JSON(404, gin.H{"error": "website not found"})
		return
	}

	c.JSON(200, website)

}
