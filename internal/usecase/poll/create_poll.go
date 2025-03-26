package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/model"
	"github.com/grozaqueen/poll/internal/utils"
	"time"
)

func (pm *PollUseCase) CreatePoll(poll model.Poll) (model.Poll, error) {
	nowInMSK := time.Now().In(utils.MskLocation)
	if poll.EndDate.Before(nowInMSK) {
		return model.Poll{}, errs.PollDateInPast
	}
	poll, err := pm.PollRepository.CreatePoll(poll)
	return poll, err
}
