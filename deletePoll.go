package main

import (
	"encoding/json"
	"net/http"
)

type DeletePollRequest struct {
	PollID string `json:"poll_id"`
	UserID string `json:"user_id"`
}

type DeletePollResponse struct {
	Status    string `json:"status"`
	PollID    string `json:"poll_id"`
	DeletedBy string `json:"deleted_by"`
}

func deletePollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req DeletePollRequest
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
		http.Error(w, "Только создатель может удалить опрос", http.StatusForbidden)
		return
	}

	delete(polls, req.PollID)

	response := DeletePollResponse{
		Status: "Опрос успешно удален",
		PollID: req.PollID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}

}
