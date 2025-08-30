package producer

import (
	"context"
	"fmt"

	"github.com/BetterStack/db"
	redisclient "github.com/BetterStack/redis"
	"github.com/redis/go-redis/v9"
)

func Producer() {
	fmt.Println("Producer is starting...")
	// Context is used for early cancellation
	ctx := context.Background()
	client := redisclient.Get()
	dbclient := db.GetPrismaClient()
	websites, err := dbclient.Website.FindMany().Exec(ctx)
	if err != nil {
		panic(err)
	}
	for _, website := range websites {
		client.XAdd(ctx, &redis.XAddArgs{
			Stream: "betterstack:website",
			Values: map[string]interface{}{
				"url": website.URL,
				"id":  website.ID},
		})
	}
}
