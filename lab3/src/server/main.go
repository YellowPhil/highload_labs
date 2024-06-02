package main

import (
	"context"
	"log"
    "os"
    "fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

const hashKey = "hash"

func GenerateClient() {
    url := fmt.Sprintf("redis://%s:%s@redis:6379/%s?protocol=3",
        os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASSWORD"),os.Getenv("REDIS_DB"))
    opts, err := redis.ParseURL(url)
    if err != nil {
        panic(err)
    }
    client := redis.NewClient(opts)

    _, err = client.Ping(context.Background()).Result()
    if err != nil {
        panic(err)
    }
    redisClient = client
}


func AddHashToBrute(hash string) {
    err := redisClient.LPush(context.Background(), hashKey, hash).Err()
    if err != nil {
        log.Println(err)
    }
}

func main() {
    GenerateClient()
    app := fiber.New()
    app.Post("/add_task", func(ctx fiber.Ctx) error {
        hashRequest := &struct {
            Hash string `json:"hash" xml:"hash" form:"hash"`
        }{}
        if err := ctx.Bind().Body(hashRequest); err != nil {
            log.Println(err)
            return err
        }
        AddHashToBrute(hashRequest.Hash)
        return nil
    })
    app.Listen(":8080")
}
