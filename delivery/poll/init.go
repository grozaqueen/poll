package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/model"
)

type PollUsecase interface {
	CreatePoll(poll model.Poll) (model.Poll, error)
}

type PollRepository interface {
	CompletePollEarly(pollID string, userID string) (model.Poll, error)
	DeletePoll(pollID string, userID string) error
	GetResults(pollID string) (model.Poll, error)
}

type PollDelivery struct {
	PollUsecase    PollUsecase
	PollRepository PollRepository
	errResolver    errs.GetErrorCode
}

func NewPollDelivery(PollUsecase PollUsecase, PollRepository PollRepository,
	errResolver errs.GetErrorCode) *PollDelivery {
	return &PollDelivery{
		PollUsecase:    PollUsecase,
		PollRepository: PollRepository,
		errResolver:    errResolver,
	}
}
