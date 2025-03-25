package poll

import (
	"encoding/json"
	"net/http"
)

func (pd *PollDelivery) GetResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	pollID := r.URL.Query().Get("id")
	if pollID == "" {
		http.Error(w, "Не указан ID опроса", http.StatusBadRequest)
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

	json.NewEncoder(w).Encode(response)
}
