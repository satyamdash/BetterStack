package worker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/BetterStack/db"
	redisclient "github.com/BetterStack/redis"
	"github.com/redis/go-redis/v9"
)

var (
	REGION_ID = "INDIA"
	WORKER_ID = "india-a"
)

func Worker() {
	fmt.Println("Worker is starting...")
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found or could not load it")
	// }
	fmt.Println("REGION_ID:", REGION_ID, "WORKER_ID:", WORKER_ID)
	ctx := context.Background()
	client := redisclient.Get()

	// Create group if not exists
	err := client.XGroupCreateMkStream(ctx, "betterstack:website", REGION_ID, "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		panic(err)
	}

	dbclient := db.GetPrismaClient()

	for {
		msgs, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    REGION_ID,
			Consumer: WORKER_ID,
			Streams:  []string{"betterstack:website", ">"},
			Block:    0,
			Count:    1,
		}).Result()

		if err != nil {
			fmt.Println("XReadGroup error:", err)
			continue
		}

		for _, stream := range msgs {
			for _, m := range stream.Messages {
				websiteID, ok1 := m.Values["id"].(string)
				websiteURL, ok2 := m.Values["url"].(string)
				if !ok1 || !ok2 {
					fmt.Println("Malformed message:", m.Values)
					continue
				}

				// Process concurrently
				go func(id, url, msgID string) {
					fetchWebsite(ctx, dbclient, client, id, url, msgID)
				}(websiteID, websiteURL, m.ID)
			}
		}
	}
}

func fetchWebsite(ctx context.Context, dbclient *db.PrismaClient, client *redis.Client, websiteID, websiteURL, msgID string) {
	fmt.Println("Fetching website:", websiteURL)
	start := time.Now()
	resp, err := http.Get(websiteURL)
	fmt.Println("Response received for:", websiteURL)
	if err != nil {
		fmt.Println("Error fetching website:", err)
		dbclient.WebsiteTick.CreateOne(
			db.WebsiteTick.ResponseTimeMs.Set(0),
			db.WebsiteTick.Status.Set(db.WebsiteStatusDown),
			db.WebsiteTick.Region.Link(db.Region.ID.Equals(REGION_ID)),
			db.WebsiteTick.Website.Link(db.Website.ID.Equals(websiteID)),
		).Exec(ctx)
		_ = client.XAck(ctx, "betterstack:website", REGION_ID, msgID).Err()
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	fmt.Printf("Website %s responded in %v with status %s\n", websiteURL, duration, resp.Status)
	res, err := dbclient.WebsiteTick.CreateOne(
		db.WebsiteTick.ResponseTimeMs.Set(int(duration.Milliseconds())),
		db.WebsiteTick.Status.Set(db.WebsiteStatusUp),
		db.WebsiteTick.Region.Link(db.Region.ID.Equals(REGION_ID)),
		db.WebsiteTick.Website.Link(db.Website.ID.Equals(websiteID)),
	).Exec(ctx)

	if err != nil {
		fmt.Println("❌ Failed to insert WebsiteTick:", err)
		return
	}

	fmt.Println("✅ Inserted WebsiteTick:", res.ID)

	// Acknowledge
	_ = client.XAck(ctx, "betterstack:website", REGION_ID, msgID).Err()
}
