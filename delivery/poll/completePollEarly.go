package poll

import (
	"encoding/json"
	"net/http"
)

func (pd *PollDelivery) CompletePollEarly(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req CompletePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	poll, err := pd.PollRepository.CompletePollEarly(req.PollID, req.UserID)
	if err != nil {
		err, statusCode := pd.errResolver.Get(err)
		http.Error(w, err.Error(), statusCode)
		return
	}
	response := CompletePollResponse{
		Status:  "Голосование завершено досрочно",
		EndDate: poll.EndDate.Format("02.01.2006 15:04:05 MST"),
		PollID:  req.PollID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}
}
