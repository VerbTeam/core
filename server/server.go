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
	main := log.New(os.Stdout, "[SERVER]: ", log.Ldate|log.Ltime|log.Lshortfile)

	main.Println("init redis...")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_PUBLIC_ENDPOINT"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORDS"),
		DB:       0,
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, req *http.Request) {
		main.Println("new request for /submit!")
		ctx := context.Background()

		value := req.URL.Query().Get("value")
		if value == "" {
			main.Println("bad request")
			sendJSONError(w, "bad request", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(value)
		if err != nil {
			main.Println("bad request")
			sendJSONError(w, "value must be an int", http.StatusBadRequest)
			return
		}

		main.Println("getting cached value...")
		cacheKey := fmt.Sprintf("worker:%d", id)
		cached, err := redisClient.Get(ctx, cacheKey).Result()
		var resp map[string]interface{}

		if err == redis.Nil {
			main.Println("nothing cached. requesting new worker")

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
		main.Print("End for /submit")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	main.Printf("server now live at https://localhost:%v\n", port)
	main.Fatal(http.ListenAndServe(":"+port, nil))

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
	main := log.New(os.Stdout, "[WORKER]: ", log.Ldate|log.Ltime|log.Lshortfile)

	main.Println("new worker running")

	type result struct{ name, data string }
	resChan := make(chan result)

	main.Println("new thread for bio")
	go func() { resChan <- result{"bio", workers.BioRun(id)} }()

	main.Println("new thread for bioAI")
	go func() { resChan <- result{"bioAI", workers.BioRunAI(id)} }()

	main.Println("new thread for avatar")
	go func() { resChan <- result{"avatar", workers.AvatarRun(id)} }()

	main.Println("new thread for group")
	go func() { resChan <- result{"group", workers.RunGroupCheck(id)} }()

	bioArray := []interface{}{}
	bioAIMap := map[string]interface{}{}
	avatarMap := map[string]interface{}{}
	groupArray := []interface{}{}

	for i := 0; i < 4; i++ {
		r := <-resChan
		switch r.name {
		case "bio":
			main.Println("rev bio")
			json.Unmarshal([]byte(r.data), &bioArray)
		case "group":
			main.Println("rev group")
			json.Unmarshal([]byte(r.data), &groupArray)
		case "bioAI":
			main.Println("rev bioAI")
			var unescaped string
			if len(r.data) > 0 && r.data[0] == '"' {
				json.Unmarshal([]byte(r.data), &unescaped)
				r.data = unescaped
			}
			err := json.Unmarshal([]byte(r.data), &bioAIMap)
			if err != nil {
				main.Println("bioAI unmarshal error:", err)
				main.Println("bioAI raw data:", r.data)
			}
		case "avatar":
			main.Println("rev AI")
			err := json.Unmarshal([]byte(r.data), &avatarMap)
			if err != nil {
				main.Println("avatar unmarshal error:", err)
				main.Println("avatar raw data:", r.data)
			}
		}
	}

	finalMap := map[string]interface{}{
		"bio": map[string]interface{}{
			"bloxdbwordlist": bioArray,
			"bioAI":          bioAIMap,
		},
		"avatar":        avatarMap,
		"flaggedGroups": groupArray,
	}

	finalJSON, _ := json.MarshalIndent(finalMap, "", "  ")

	main.Println("done")

	return string(finalJSON)
}
