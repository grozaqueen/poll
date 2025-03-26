package poll

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (pd *PollDelivery) GetResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	pollIDStr := r.URL.Query().Get("id")
	if pollIDStr == "" {
		http.Error(w, "Не указан ID опроса", http.StatusBadRequest)
		return
	}

	// Конвертируем строку в uint64
	pollID, err := strconv.ParseUint(pollIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID опроса", http.StatusBadRequest)
		return
	}

	poll, err := pd.PollRepository.GetResults(pollID)
	if err != nil {
		err, statusCode := pd.errResolver.Get(err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	response := ResultsPollResponse{
		Question: poll.Question,
		Options:  poll.Options,
		Votes:    poll.Votes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
