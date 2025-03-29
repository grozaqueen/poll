package poll

import (
	"net/http"
)

const formatDate = "02.01.2006 15:04:05 MST"

// CompletePollEarly godoc
// @Summary Завершить опрос досрочно
// @Description Позволяет создателю завершить опрос до установленного срока окончания
// @Tags Polls
// @Accept json
// @Produce json
// @Param input body CompletePollRequest true "Данные для досрочного завершения опроса"
// @Success 200 {object} CompletePollResponse "Опрос успешно завершен"
// @Failure 400 {object} string "Неверный формат JSON"
// @Failure 403 {object} string "Только создатель может выполнить это действие"
// @Failure 404 {object} string "Опрос не найден"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /poll/complete [patch]
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
