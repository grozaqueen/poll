package vote

import (
	"net/http"
)

// CreateVote godoc
// @Summary Создать голос
// @Description Записывает голос пользователя в указанном опросе
// @Tags Votes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body CreateVoteRequest true "Данные голоса"
// @Success 200 {object} string "Голос успешно учтен"
// @Failure 400 {object} string "Неверный формат запроса"
// @Failure 403 {object} string "Вы уже голосовали в этом опросе"
// @Failure 404 {object} string "Опрос не найден"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /vote [post]
func (vd *VoteDelivery) CreateVote(w http.ResponseWriter, r *http.Request) {
	if !vd.utils.ValidateMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateVoteRequest
	if !vd.utils.DecodeRequest(w, r, &req) {
		return
	}

	err := vd.VoteCreator.CreateVote(req.toModel())
	if err != nil {
		vd.utils.HandleError(w, r, err, "CreateVote: ошибка создания голоса")
		return
	}

	vd.utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status": "Голос учтён!",
	})
}
