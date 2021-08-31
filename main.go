package main

import (
	"log"
	"net/http"
	"os"
	. "self/prelude"
	"self/updating"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
)

func main() {
	redisUrl, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		redisUrl = "localhost:6379"
	}

	// conn redis
	log.Println("connecting redis:", redisUrl)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// run thread that handles updating data
	go updating.Run(rdb)

	// start http server
	Serve(rdb)
}

// Serve creates an http server and listens.
func Serve(rdb *redis.Client) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/chat-commands/twitch/{uid}", func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, "uid")

		updating.HintUpdate(uid)
		val, err := rdb.Get(RedisCtx, GetRedisKey4("twitch", uid, "commands", "all")).Result()

		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(val))
	})

	log.Println("listening at http://localhost:9965")
	log.Println("http server exit:", http.ListenAndServe(":9965", r))
}
