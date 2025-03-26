package vote

import (
	"encoding/json"
	"net/http"
)

func (vd *VoteDelivery) CreateVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req CreateVoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	err := vd.VoteCreator.CreateVote(req.toModel())
	if err != nil {
		err, statusCode := vd.errResolver.Get(err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Голос учтён!",
	})
}
