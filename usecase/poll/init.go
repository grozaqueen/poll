package poll

import (
	"github.com/grozaqueen/poll/model"
)

type PollRepository interface {
	CreatePoll(poll model.Poll) (model.Poll, error)
}

type PollUseCase struct {
	PollRepository PollRepository
}

func NewPollUseCase(pollRepository PollRepository) *PollUseCase {
	return &PollUseCase{
		PollRepository: pollRepository,
	}
}
