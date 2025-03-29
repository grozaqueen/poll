package poll

import (
	"net/http"
)

// DeletePoll godoc
// @Summary Удалить опрос
// @Description Удаляет существующий опрос (доступно только создателю опроса)
// @Tags Polls
// @Accept json
// @Produce json
// @Param input body DeletePollRequest true "Данные для удаления опроса"
// @Success 200 {object} DeletePollResponse "Опрос успешно удален"
// @Failure 400 {object} string "Неверный формат JSON"
// @Failure 403 {object} string "Только создатель может выполнить это действие"
// @Failure 404 {object} string "Опрос не найден"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /poll/delete [delete]
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
