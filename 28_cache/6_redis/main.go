package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("redis close failed: %v", err)
		}
	}()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	fmt.Println("connected to redis")

	cacheKey := "user:1"

	if err := client.Set(ctx, cacheKey, "Алексей", 5*time.Second).Err(); err != nil {
		log.Fatalf("redis set failed: %v", err)
	}
	fmt.Println("value cached with TTL 5s")

	val, err := client.Get(ctx, cacheKey).Result()
	if err != nil {
		log.Fatalf("redis get failed: %v", err)
	}
	fmt.Printf("cached value: %s\n", val)

	fmt.Println("sleeping 6s to show expiration...")
	time.Sleep(6 * time.Second)

	_, err = client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		fmt.Println("value expired and is absent in cache")
	} else if err != nil {
		log.Fatalf("redis get after ttl failed: %v", err)
	}
}
