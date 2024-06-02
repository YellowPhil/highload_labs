package main

import (
	"context"
	"fmt"
	"highload/worker/bruter"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const hashKey = "hash"

var ctx = context.Background()
var redisClient *redis.Client

func GenerateClient() {
    url := fmt.Sprintf("redis://%s:%s@redis:6379/%s?protocol=3",
        os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASSWORD"),os.Getenv("REDIS_DB"))
    opts, err := redis.ParseURL(url)
    if err != nil {
        panic(err)
    }
    client := redis.NewClient(opts)

    _, err = client.Ping(ctx).Result()
    if err != nil {
        panic(err)
    }
    redisClient = client
}

func worker(hash string) bruter.BruteforceResult {
    return bruter.BruteHash(hash, bruter.SHA256)
}

func main() {
    bruter.Prepare(os.Getenv("WORDLIST_PATH"))
    GenerateClient()

    wp := bruter.NewWorkerPool[string, bruter.BruteforceResult](
        16, time.Minute * 3, 
        worker, bruter.IsError,)
    wp.Start(
        func (input string) { log.Printf("[-] Unable to recover password for hash: %v\n", input)},
        func(result bruter.BruteforceResult) { log.Printf("[+] Recovered hash: %v : %v\n", result.Input, string(result.Recovered))},
        )

    for {
        ctx := context.Background()
        if val, err := redisClient.BRPop(ctx, 0, hashKey).Result(); err != nil {
            log.Fatal("ERROR: no such key")
        } else {
            log.Printf("[!] Got value %v\n", val)
            wp.AddTask(val[1])
        }
    }
    GenerateClient()
}
