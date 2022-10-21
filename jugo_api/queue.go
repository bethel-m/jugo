package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
)

func init() {
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_db, conv_err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if conv_err != nil {
		log.Fatalf("could not convert redis db to number::%v", conv_err)
	}

	redis_address := fmt.Sprintf("%v:%v", redis_host, redis_port)
	ctx := context.Background()
	client = redis.NewClient(&redis.Options{
		Addr:     redis_address,
		Password: redis_password,
		DB:       redis_db,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("error pinging redis::%v", err)
	}
	fmt.Println(pong)
	fmt.Println("this is another one")
}

func add_task_to_queue(client *redis.Client, username string) (int64, error) {
	ctx := context.Background()
	task_queue := client.LPush(ctx, "tasks_queue", username)

	no_of_tasks_in_queue := task_queue.Val()
	task_queue_err := task_queue.Err()
	if task_queue_err != nil {
		fmt.Printf("error adding task to queue::%v\n", username)
		return no_of_tasks_in_queue, task_queue_err
	}
	fmt.Printf("task for %v, added to queue\n", username)
	return no_of_tasks_in_queue, nil
}
