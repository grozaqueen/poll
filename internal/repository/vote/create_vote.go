package vote

import (
	"errors"
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/model"
	"strconv"
)

func (vr *VoteRepository) CreateVote(vote model.Vote) error {
	userIDStr := strconv.FormatUint(vote.UserID, 10)

	resp, err := vr.polRep.Call("create_vote", []interface{}{vote.PollID, vote.Option, userIDStr})
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return errors.New("empty response from Tarantool")
	}

	result, ok := resp[0].(map[interface{}]interface{})
	if !ok {
		return errors.New("invalid response format")
	}

	if errMsg, exists := result["error"]; exists {
		switch errMsg.(string) {
		case "PollNotFound":
			return errs.PollNotFound
		case "PollAlreadyClosed":
			return errs.PollAlreadyClosed
		case "UserAlreadyVoted":
			return errs.UserAlreadyVoted
		case "InvalidVoteOption":
			return errs.InvalidVoteOption
		default:
			return errors.New(errMsg.(string))
		}
	}

	return nil
}
