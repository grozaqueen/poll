package poll

import (
	"net/http"
)

func (pd *PollDelivery) GetResults(w http.ResponseWriter, r *http.Request) {
	if !pd.utils.ValidateMethod(w, r, http.MethodGet) {
		return
	}

	pollID, ok := pd.utils.ParseUintParam(w, r, r.URL.Query().Get("id"))
	if !ok {
		return
	}

	poll, err := pd.PollRepository.GetResults(pollID)
	if err != nil {
		pd.utils.HandleError(w, r, err, "GetResults: ошибка получения результатов")
		return
	}

	response := ResultsPollResponse{
		Question: poll.Question,
		Options:  poll.Options,
		Votes:    poll.Votes,
	}

	pd.utils.SendJSONResponse(w, http.StatusOK, response)
}
