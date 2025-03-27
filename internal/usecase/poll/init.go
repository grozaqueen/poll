package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/model"
	"log/slog"
)

type PollRepository interface {
	CreatePoll(poll model.Poll) (model.Poll, error)
}

type PollUseCase struct {
	PollRepository PollRepository
	log            *slog.Logger
	errResolver    errs.GetErrorCode
}

func NewPollUseCase(pollRepository PollRepository, log *slog.Logger,
	errResolver errs.GetErrorCode) *PollUseCase {
	return &PollUseCase{
		PollRepository: pollRepository,
		log:            log,
		errResolver:    errResolver,
	}
}
