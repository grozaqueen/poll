package vote

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/model"
	"time"
)

func (vr *VoteRepository) CreateVote(vote model.Vote) error {
	vr.polRep.PollsMu.RLock()
	poll, exists := vr.polRep.Polls[vote.PollID]
	vr.polRep.PollsMu.RUnlock()

	if !exists {
		return errs.PollNotFound
	}

	if time.Now().After(poll.EndDate) {
		return errs.PollAlreadyClosed
	}

	vr.polRep.PollsMu.Lock()
	defer vr.polRep.PollsMu.Unlock()

	if poll.Voters == nil {
		poll.Voters = make(map[string]bool)
	}

	if _, voted := poll.Voters[vote.UserID]; voted {
		return errs.UserAlreadyVoted
	}

	validOption := false
	for _, opt := range poll.Options {
		if opt == vote.Option {
			validOption = true
			break
		}
	}
	if !validOption {
		return errs.InvalidVoteOption
	}

	poll.Votes[vote.Option]++
	poll.Voters[vote.UserID] = true

	return nil
}
