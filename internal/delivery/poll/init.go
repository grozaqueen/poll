package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/model"
	"time"
)

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
