package vote

import (
	"net/http"
)

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
