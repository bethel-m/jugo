package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

func init() {
	ctx := context.TODO()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis_db:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.RPush(ctx, "trial", "show_up").Err()
	if err != nil {
		fmt.Println("error occures::,", err)
	} else {
		fmt.Println("sent successfully")
	}
}
