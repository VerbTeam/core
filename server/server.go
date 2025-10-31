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

func newWorker(id int) string {
	fmt.Println("new worker running")

	type result struct{ name, data string }
	resChan := make(chan result)

	go func() { resChan <- result{"bio", workers.BioRun(id)} }()
	go func() { resChan <- result{"bioAI", workers.BioRunAI(id)} }()
	go func() { resChan <- result{"avatar", workers.AvatarRun(id)} }()

	bioArray := []interface{}{}
	bioAIMap := map[string]interface{}{}
	avatarMap := map[string]interface{}{}

	for i := 0; i < 3; i++ {
		r := <-resChan
		switch r.name {
		case "bio":
			json.Unmarshal([]byte(r.data), &bioArray)
		case "bioAI":
			var unescaped string
			if len(r.data) > 0 && r.data[0] == '"' {
				json.Unmarshal([]byte(r.data), &unescaped)
				r.data = unescaped
			}
			err := json.Unmarshal([]byte(r.data), &bioAIMap)
			if err != nil {
				fmt.Println("bioAI unmarshal error:", err)
				fmt.Println("bioAI raw data:", r.data)
			}
		case "avatar":
			err := json.Unmarshal([]byte(r.data), &avatarMap)
			if err != nil {
				fmt.Println("avatar unmarshal error:", err)
				fmt.Println("avatar raw data:", r.data)
			}
		}
	}

	finalMap := map[string]interface{}{
		"bio": map[string]interface{}{
			"bloxdbwordlist": bioArray,
			"bioAI":          bioAIMap,
		},
		"avatar": avatarMap,
	}

	finalJSON, _ := json.MarshalIndent(finalMap, "", "  ")
	return string(finalJSON)
}
