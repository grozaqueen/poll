package main

import (
	"encoding/json"
	"net/http"
)

func getResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	pollID := r.URL.Query().Get("id")
	if pollID == "" {
		http.Error(w, "Не указан ID опроса", http.StatusBadRequest)
		return
	}

	pollsMu.RLock()
	poll, exists := polls[pollID]
	pollsMu.RUnlock()

	if !exists {
		http.Error(w, "Опрос не найден", http.StatusNotFound)
		return
	}

	response := struct {
		Question string         `json:"question"`
		Options  []string       `json:"options"`
		Votes    map[string]int `json:"votes"`
	}{
		Question: poll.Question,
		Options:  poll.Options,
		Votes:    poll.Votes,
	}

	json.NewEncoder(w).Encode(response)
}
