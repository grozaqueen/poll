package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type CompletePollRequest struct {
	PollID string `json:"poll_id"`
	UserID string `json:"user_id"`
}

type CompletePollResponse struct {
	Status      string `json:"status"`
	EndDate     string `json:"end_date"`
	PollID      string `json:"poll_id"`
	OriginalEnd string `json:"original_end,omitempty"`
}

func completePollEarly(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req CompletePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	pollsMu.Lock()
	defer pollsMu.Unlock()

	poll, exists := polls[req.PollID]
	if !exists {
		http.Error(w, "Опрос не найден", http.StatusNotFound)
		return
	}

	if poll.Creator.ID != req.UserID {
		http.Error(w, "Только создатель может завершить опрос", http.StatusForbidden)
		return
	}

	originalEnd := poll.EndDate

	poll.EndDate = time.Now().In(mskLocation)

	response := CompletePollResponse{
		Status:      "Голосование завершено досрочно",
		EndDate:     poll.EndDate.Format("02.01.2006 15:04:05 MST"),
		PollID:      req.PollID,
		OriginalEnd: originalEnd.Format("02.01.2006 15:04:05 MST"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}
}
