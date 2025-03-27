package poll

import (
	"net/http"
)

func (pd *PollDelivery) DeletePoll(w http.ResponseWriter, r *http.Request) {
	if !pd.utils.ValidateMethod(w, r, http.MethodDelete) {
		return
	}

	var req DeletePollRequest
	if !pd.utils.DecodeRequest(w, r, &req) {
		return
	}

	err := pd.PollRepository.DeletePoll(req.PollID, req.UserID)
	if err != nil {
		pd.utils.HandleError(w, r, err, "DeletePoll: ошибка удаления опроса")
		return
	}

	response := DeletePollResponse{
		Status: "Опрос успешно удален",
		PollID: req.PollID,
	}

	pd.utils.SendJSONResponse(w, http.StatusOK, response)
}
