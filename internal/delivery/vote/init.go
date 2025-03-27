package vote

import (
	"github.com/grozaqueen/poll/internal/model"
	"net/http"
)

type UtilsHelpers interface {
	ValidateMethod(w http.ResponseWriter, r *http.Request, method string) bool
	DecodeRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool
	SendJSONResponse(w http.ResponseWriter, status int, data interface{})
	HandleError(w http.ResponseWriter, r *http.Request, err error, context string)
	ParseUintParam(w http.ResponseWriter, r *http.Request, param string) (uint64, bool)
}

type VoteCreator interface {
	CreateVote(vote model.Vote) error
}

type VoteDelivery struct {
	VoteCreator VoteCreator
	utils       UtilsHelpers
}

func NewVoteDelivery(VoteCreator VoteCreator, utils UtilsHelpers) *VoteDelivery {
	return &VoteDelivery{
		VoteCreator: VoteCreator,
		utils:       utils,
	}
}
