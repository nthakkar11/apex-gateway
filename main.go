package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
	// Atomic Lua script for Idempotency and Rate Limiting
	gatekeeperScript = redis.NewScript(`
		local user_limit_key = KEYS[1]
		local idempotency_key = KEYS[2]
		local limit = tonumber(ARGV[1])

		local cached_res = redis.call("GET", idempotency_key)
		if cached_res then return {200, cached_res} end

		local current = redis.call("INCR", user_limit_key)
		if current == 1 then redis.call("EXPIRE", user_limit_key, 60) end

		if current > limit then return {429, "Rate limit exceeded"} end

		return {201, "OK"}
	`)
)

func main() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("Invalid Redis URL")
	}

	opt.PoolSize = 100
	opt.MinIdleConns = 10
	opt.DialTimeout = 5 * time.Second

	rdb = redis.NewClient(opt)

	http.HandleFunc("/v1/transaction", handleTransaction)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Gateway operational on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Idempotency-Key")
	uID := r.URL.Query().Get("user_id")
	iKey := r.Header.Get("X-Idempotency-Key")

	if uID == "" || iKey == "" {
		http.Error(w, "Incomplete request parameters", http.StatusBadRequest)
		return
	}

	res, err := gatekeeperScript.Run(ctx, rdb, []string{"lim:" + uID, "idem:" + iKey}, 100).Slice()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	status := int(res[0].(int64))
	msg := res[1].(string)

	if status == 201 {
		// Simulate logic execution and cache result
		txID := fmt.Sprintf("TXN_%d", time.Now().UnixNano())
		rdb.Set(ctx, "idem:"+iKey, txID, 24*time.Hour)
		w.WriteHeader(http.StatusCreated)
		log.Printf("[INFO] SUCCESS: Created Transaction %s for User %s", txID, uID)
		return
	} else if status == 200 {
		log.Printf("[INFO] IDEMPOTENCY HIT: Returning cached result for Key %s", iKey)
	} else if status == 429 {
		log.Printf("[WARN] RATE LIMIT: User %s exceeded threshold", uID)
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, `{"message": "%s"}`, msg)
}
