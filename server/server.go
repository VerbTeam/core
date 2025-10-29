package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	workers "codeberg.org/VerbTeam/core/worker"
	"github.com/redis/go-redis/v9"
)

func Start() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_PUBLIC_ENDPOINT"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORDS"),
		DB:       0,
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()

		value := req.URL.Query().Get("value")
		if value == "" {
			sendJSONError(w, "bad request", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(value)
		if err != nil {
			sendJSONError(w, "value must be an int", http.StatusBadRequest)
			return
		}

		cacheKey := fmt.Sprintf("worker:%d", id)
		cached, err := redisClient.Get(ctx, cacheKey).Result()
		var resp map[string]interface{}

		if err == redis.Nil {
			res := newWorker(id)
			var parsed map[string]interface{}
			json.Unmarshal([]byte(res), &parsed)
			resp = parsed

			jsonData, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}

			err = redisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err()
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		} else {
			json.Unmarshal([]byte(cached), &resp)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("server now live at https://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func sendJSONError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errResp := map[string]string{
		"status":  "error",
		"message": msg,
	}
	json.NewEncoder(w).Encode(errResp)
}

type Avatar struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type ProcessedInformation struct {
	BioReason []Avatar `json:"bio"`
	Avatar    Avatar   `json:"avatar"`
}

func newWorker(id int) string {
	fmt.Println("new worker running")

	type result struct {
		name string
		data string
	}

	resChan := make(chan result)

	go func() { resChan <- result{"bio", workers.BioRun(id)} }()
	go func() { resChan <- result{"avatar", workers.AvatarRun(id)} }()

	var bioData, avatarData string
	for i := 0; i < 2; i++ {
		r := <-resChan
		if r.name == "bio" {
			bioData = r.data
		} else if r.name == "avatar" {
			avatarData = r.data
		}
	}

	final := fmt.Sprintf(`{"bio":%s,"avatar":%s}`, bioData, avatarData)
	return final
}
