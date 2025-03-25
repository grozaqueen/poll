package poll

import (
	"github.com/grozaqueen/poll/errs"
)

func (pr *PollRepository) DeletePoll(pollID string, userID string) error {
	pr.PollsMu.Lock()
	defer pr.PollsMu.Unlock()

	poll, exists := pr.Polls[pollID]
	if !exists {
		return errs.PollNotFound
	}

	if poll.Creator.ID != userID {
		return errs.UserNotCreator
	}

	delete(pr.Polls, pollID)
	return nil
}
