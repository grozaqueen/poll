package poll

import (
	"net/http"
)

// GetResults godoc
// @Summary Получить результаты опроса
// @Description Возвращает текущие результаты голосования по указанному опросу
// @Tags Polls
// @Accept json
// @Produce json
// @Param id query integer true "ID опроса" Example(123)
// @Success 200 {object} ResultsPollResponse "Успешный запрос"
// @Failure 400 {object} string "Неверный формат ID"
// @Failure 404 {object} string "Опрос не найден"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /results [get]
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
