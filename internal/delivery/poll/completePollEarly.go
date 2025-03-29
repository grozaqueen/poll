package poll

import (
	"net/http"
)

const formatDate = "02.01.2006 15:04:05 MST"

func (pd *PollDelivery) CompletePollEarly(w http.ResponseWriter, r *http.Request) {
	if !pd.utils.ValidateMethod(w, r, http.MethodPatch) {
		return
	}

	var req CompletePollRequest
	if !pd.utils.DecodeRequest(w, r, &req) {
		return
	}

	endDate, err := pd.PollRepository.CompletePollEarly(req.PollID, req.UserID)
	if err != nil {
		pd.utils.HandleError(w, r, err, "CompletePollEarly: ошибка завершения опроса")
		return
	}

	response := CompletePollResponse{
		Status:  "Голосование завершено досрочно",
		EndDate: endDate.Format(formatDate),
		PollID:  req.PollID,
	}

	pd.utils.SendJSONResponse(w, http.StatusOK, response)
}
