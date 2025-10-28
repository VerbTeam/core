package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"strconv"

	"github.com/redis/go-redis/v9"

	workers "codeberg.org/VerbTeam/core/worker/"
)

var ctx = context.Background()

func Start() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_PUBLIC_ENDPOINT"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORDS"),
		DB:       0,
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")

		if value == "" {
			sendJSONError(w, "bad request", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(value)
		if err != nil {
			sendJSONError(w, "value must be a int", http.StatusBadRequest)
			return
		}

		resp := map[string]string{
			"status": "submited!",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

		go newWorker(id)
	})

	fmt.Println("server up at http://localhost:8080")
	http.ListenAndServe(":8080", nil)

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

func newWorker(id int) {

	resChan := make(chan string) // make a channel

	// run this in a diffrent thread (or that what i know)
	go func() {
		resChan <- workers.BioRun(id) // catch it
	}()

	go func() {
		resChan <- workers.AvatarRun(id) // catch it
	}()

	resBio := <-resChan // goes to this
	resAvatar := <-resChan

	fmt.Println(resBio)
	fmt.Println(resAvatar)
}
