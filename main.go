package main

import (
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
)

var (
	pool             *redis.Pool
	redisServer, _   = os.LookupEnv("REDIS_URL")
	redisPassword, _ = os.LookupEnv("REDIS_PASSWORD")
)

func main() {
	pool = newPool(redisServer, redisPassword)
	defer pool.Close()

	loadEffects()

	log.Printf("Starting server")

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("./public/")).ServeHTTP(w, r)
	})

	router.HandleFunc("/magical-effect", randomMagicalEffect)

	log.Fatal(http.ListenAndServe(":8080", router))
}
