package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/model"
)

func (pr *PollRepository) GetResults(pollID string) (model.Poll, error) {
	pr.PollsMu.RLock()
	poll, exists := pr.Polls[pollID]
	pr.PollsMu.RUnlock()

	if !exists {
		return model.Poll{}, errs.PollNotFound
	}

	return *poll, nil
}
