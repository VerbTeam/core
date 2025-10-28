package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	workers "codeberg.org/VerbTeam/core/worker"
)

func Start() {

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

		res := newWorker(id)

		resp := map[string]string{
			"status": res,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

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

type Avatar struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type ProcessedInformation struct {
	BioReason []Avatar `json:"bio"`
	Avatar    Avatar   `json:"avatar"`
}

func newWorker(id int) string {
	fmt.Println("starting new worker..")

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

	var bio []Avatar
	if err := json.Unmarshal([]byte(bioData), &bio); err != nil {
		panic(err)
	}

	var avatar Avatar
	if err := json.Unmarshal([]byte(avatarData), &avatar); err != nil {
		panic(err)
	}

	processed := ProcessedInformation{
		BioReason: bio,
		Avatar:    avatar,
	}

	jsonBytes, err := json.Marshal(processed)
	if err != nil {
		panic(err)
	}

	return string(jsonBytes)
}
