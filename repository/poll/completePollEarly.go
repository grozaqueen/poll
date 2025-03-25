package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/model"
	"github.com/grozaqueen/poll/utils"
	"time"
)

func (pr *PollRepository) CompletePollEarly(pollID string, userID string) (model.Poll, error) {
	pr.PollsMu.Lock()
	defer pr.PollsMu.Unlock()

	poll, exists := pr.Polls[pollID]
	if !exists {
		return model.Poll{}, errs.PollNotFound
	}

	if poll.Creator.ID != userID {
		return model.Poll{}, errs.UserNotCreator
	}

	poll.EndDate = time.Now().In(utils.MskLocation)

	return *poll, nil
}
