package poll

import (
	"github.com/grozaqueen/poll/internal/utils"
	"net/http"
)

func (pd *PollDelivery) CreatePoll(w http.ResponseWriter, r *http.Request) {
	if !pd.utils.ValidateMethod(w, r, http.MethodPost) {
		return
	}

	var req CreatePollRequest
	if !pd.utils.DecodeRequest(w, r, &req) {
		return
	}

	endDate, err := utils.ParseSimpleDate(req.EndDate)
	if err != nil {
		pd.utils.HandleError(w, r, err, "CreatePoll: ошибка формата даты")
		return
	}

	poll, err := pd.PollUsecase.CreatePoll(req.toModel(endDate))
	if err != nil {
		pd.utils.HandleError(w, r, err, "CreatePoll: ошибка создания опроса")
		return
	}

	response := CreatePollResponse{
		PollID:   poll.ID,
		Options:  poll.Options,
		UserID:   poll.Creator.ID,
		UserName: poll.Creator.Name,
	}

	pd.utils.SendJSONResponse(w, http.StatusCreated, response)
}
