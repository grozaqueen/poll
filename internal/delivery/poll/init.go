package poll

import (
	"github.com/grozaqueen/poll/internal/model"
	"net/http"
	"time"
)

type UtilsHelpers interface {
	ValidateMethod(w http.ResponseWriter, r *http.Request, method string) bool
	DecodeRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool
	SendJSONResponse(w http.ResponseWriter, status int, data interface{})
	HandleError(w http.ResponseWriter, r *http.Request, err error, context string)
	ParseUintParam(w http.ResponseWriter, r *http.Request, param string) (uint64, bool)
}

type PollUsecase interface {
	CreatePoll(poll model.Poll) (model.Poll, error)
}

type PollRepository interface {
	CompletePollEarly(pollID uint64, userID uint64) (time.Time, error)
	DeletePoll(pollID uint64, userID uint64) error
	GetResults(pollID uint64) (model.Poll, error)
}

type PollDelivery struct {
	PollUsecase    PollUsecase
	PollRepository PollRepository
	utils          UtilsHelpers
}

func NewPollDelivery(PollUsecase PollUsecase, PollRepository PollRepository,
	utils UtilsHelpers) *PollDelivery {
	return &PollDelivery{
		PollUsecase:    PollUsecase,
		PollRepository: PollRepository,
		utils:          utils,
	}
}
