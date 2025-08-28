package db

import (
	"log"
	"sync"
)

var clientInstance *PrismaClient
var once sync.Once

// GetPrismaClient returns a singleton Prisma client
func GetPrismaClient() *PrismaClient {
	once.Do(func() {
		client := NewClient()

		// connect once
		if err := client.Prisma.Connect(); err != nil {
			log.Fatalf("failed to connect to prisma: %v", err)
		}

		clientInstance = client
	})
	return clientInstance
}

// Disconnect cleans up the connection (call on shutdown)
func Disconnect() {
	if clientInstance != nil {
		if err := clientInstance.Prisma.Disconnect(); err != nil {
			log.Printf("failed to disconnect prisma: %v", err)
		}
	}
}
