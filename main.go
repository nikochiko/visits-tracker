package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

const visitsKey = "visits"

var redisAddr = getRedisAddr()
var laddr = getListenAddr()

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	l := Listener{rdb: rdb}

	http.HandleFunc("/visits", l.HandleVisits)

	log.Printf("Info: starting to serve on addr %s\n", laddr)
	log.Fatal(http.ListenAndServe(laddr, nil))
}

type Listener struct {
	rdb *redis.Client
}

func (l *Listener) HandleVisits(w http.ResponseWriter, r *http.Request) {
	if err := l.setOrIncrementCounter(visitsKey); err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (l *Listener) setOrIncrementCounter(key string) error {
	if err := l.rdb.Get(ctx, key).Err(); err != nil {
		if err == redis.Nil {
			// try to set the key
			if err = l.rdb.Set(ctx, key, 0, 0).Err(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := l.rdb.Incr(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func getRedisAddr() string {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return redisAddr
}

func getListenAddr() string {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8002"
	}

	return listenAddr
}
