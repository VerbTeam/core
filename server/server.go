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

	"github.com/redis/go-redis/v9"

	workers "codeberg.org/VerbTeam/core/worker"
)

func Start() {
	main := log.New(os.Stdout, "[SERVER]: ", log.Ldate|log.Ltime|log.Lshortfile)

	main.Println("initing redis...")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_PUBLIC_ENDPOINT"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORDS"),
		DB:       0,
	})

	http.HandleFunc("/MLchecking", func(w http.ResponseWriter, r *http.Request) {
		main.Println("new request for /MLchecking!")
		ctx := context.Background()

		value := r.URL.Query().Get("id")
		cache := r.URL.Query().Get("cache")

		if value == "" || cache == "" {
			main.Println("bad request")
			sendJSONError(w, "bad request", http.StatusBadRequest)
			return
		}

		id, erre := strconv.Atoi(value)
		if erre != nil {
			main.Println("bad request")
			sendJSONError(w, "id must be an int", http.StatusBadRequest)
			return
		}

		useCache := cache != "false"
		var resp map[string]interface{}
		cacheKey := fmt.Sprintf("worker:%d", id)

		if useCache {
			main.Println("getting cached value...")
			cached, err := redisClient.Get(ctx, cacheKey).Result()
			if err == redis.Nil {
				main.Println("nothing cached. requesting new worker")

				res := newWorkerML(id)
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
		} else {
			main.Println("cache disabled, calling new worker")

			res := newWorkerML(id)
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

		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		main.Print("End for /MLchecking")

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := `what is this diddy blud doing with the api, is blud terry?`

		w.Write([]byte(resp))
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, req *http.Request) {
		main.Println("new request for /submit!")
		ctx := context.Background()

		value := req.URL.Query().Get("userid")
		cache := req.URL.Query().Get("cache")

		if value == "" || cache == "" {
			main.Println("bad request")
			sendJSONError(w, "bad request", http.StatusBadRequest)
			return
		}

		id, erre := strconv.Atoi(value)
		if erre != nil {
			main.Println("bad request")
			sendJSONError(w, "value must be an int", http.StatusBadRequest)
			return
		}

		useCache := cache != "false"
		var resp map[string]interface{}
		cacheKey := fmt.Sprintf("worker:%d", id)

		if useCache {
			main.Println("getting cached value...")
			cached, err := redisClient.Get(ctx, cacheKey).Result()
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
		} else {
			main.Println("cache disabled, calling new worker")

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

		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		main.Print("End for /submit")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	main.Printf("server now live at http://localhost:%v\n", port)
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

	rating := float64(len(groupArray)) * 0.5

	finalMap := map[string]interface{}{
		"bio": map[string]interface{}{
			"bloxdbwordlist": bioArray,
			"bioAI":          bioAIMap,
		},
		"avatar":        avatarMap,
		"flaggedGroups": groupArray,
		"groupRating":   rating,
	}

	finalJSON, _ := json.MarshalIndent(finalMap, "", "  ")

	main.Println("done")

	return string(finalJSON)
}

func newWorkerML(id int) string {
	main := log.New(os.Stdout, "[WORKER]: ", log.Ldate|log.Ltime|log.Lshortfile)

	main.Println("new worker running")

	type result struct{ name, data string }
	resChan := make(chan result)

	main.Println("new thread for bio")
	go func() { resChan <- result{"bio", workers.BioRun(id)} }()

	main.Println("new thread for bioAI (sybauML version)")
	go func() { resChan <- result{"bioAI", workers.BioRunAIML(id)} }()

	main.Println("new thread for group")
	go func() { resChan <- result{"group", workers.RunGroupCheck(id)} }()

	bioArray := []interface{}{}
	var bioAILabel int
	groupArray := []interface{}{}

	for i := 0; i < 3; i++ {
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

			label, err := strconv.Atoi(r.data)
			if err != nil {
				main.Println("bioAI label parsing error:", err)
				bioAILabel = -1
			} else {
				bioAILabel = label
			}
		}
	}

	rating := float64(len(groupArray)) * 0.5

	finalMap := map[string]interface{}{
		"bio": map[string]interface{}{
			"bloxdbwordlist": bioArray,
			"bioAILabel":     bioAILabel,
		},
		"flaggedGroups": groupArray,
		"groupRating":   rating,
	}

	finalJSON, _ := json.MarshalIndent(finalMap, "", "  ")

	main.Println("done")

	return string(finalJSON)
}
