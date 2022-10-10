package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

func init() {

	ctx := context.Background()
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
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
