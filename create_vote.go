package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func createVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PollID string `json:"poll_id"`
		Option string `json:"option"`
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	pollsMu.RLock()
	poll, exists := polls[req.PollID]
	pollsMu.RUnlock()

	if !exists {
		http.Error(w, "Опрос не найден", http.StatusNotFound)
		return
	}

	if time.Now().After(poll.EndDate) {
		http.Error(w, "Голосование завершено", http.StatusForbidden)
		return
	}

	pollsMu.Lock()
	defer pollsMu.Unlock()

	if poll.Voters == nil {
		poll.Voters = make(map[string]bool)
	}

	if _, voted := poll.Voters[req.UserID]; voted {
		http.Error(w, "Вы уже голосовали в этом опросе", http.StatusForbidden)
		return
	}

	validOption := false
	for _, opt := range poll.Options {
		if opt == req.Option {
			validOption = true
			break
		}
	}
	if !validOption {
		http.Error(w, "Неверный вариант", http.StatusBadRequest)
		return
	}

	poll.Votes[req.Option]++
	poll.Voters[req.UserID] = true

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Голос учтён!",
	})
}
