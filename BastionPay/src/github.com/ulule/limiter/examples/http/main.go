package main

import (
	"fmt"
	"log"
	"net/http"

	redis "github.com/go-redis/redis"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	sredis "github.com/ulule/limiter/drivers/store/redis"
)

func main() {

	// Define a limit rate to 4 requests per hour.
	rate, err := limiter.NewRateFromFormatted("4-H")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a redis client.
	option, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		log.Fatal(err)
		return
	}
	client := redis.NewClient(option)

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_http_example",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a new middleware with the limiter instance.
	middleware := stdlib.NewMiddleware(limiter.New(store, rate), stdlib.WithForwardHeader(true))

	// Launch a simple server.
	http.Handle("/", middleware.Handler(http.HandlerFunc(index)))
	fmt.Println("Server is running on port 7777...")
	log.Fatal(http.ListenAndServe(":7777", nil))

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"message": "ok"}`))
}
