package poll

import (
	"encoding/json"
	"github.com/grozaqueen/poll/internal/utils"
	"net/http"
)

func (pd *PollDelivery) CreatePoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	endDate, err := utils.ParseSimpleDate(req.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	poll, err := pd.PollUsecase.CreatePoll(req.toModel(endDate))
	if err != nil {
		err, statusCode := pd.errResolver.Get(err)
		http.Error(w, err.Error(), statusCode)
		return
	}
	response := CreatePollResponse{
		PollID:   poll.ID,
		Options:  poll.Options,
		UserID:   poll.Creator.ID,
		UserName: poll.Creator.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
