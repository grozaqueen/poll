package poll

import (
	"encoding/json"
	"net/http"
)

func (pd *PollDelivery) DeletePoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req DeletePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	err := pd.PollRepository.DeletePoll(req.PollID, req.UserID)
	if err != nil {
		err, statusCode := pd.errResolver.Get(err)
		http.Error(w, err.Error(), statusCode)
		return
	}
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
