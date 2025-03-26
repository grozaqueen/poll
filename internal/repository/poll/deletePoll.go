package poll

import (
	"errors"
	"github.com/grozaqueen/poll/errs"
)

func (pr *PollRepository) DeletePoll(pollID uint64, userID uint64) error {
	resp, err := pr.Call("delete_poll", []interface{}{pollID, userID})
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
		case "UserNotCreator":
			return errs.UserNotCreator
		default:
			return errors.New(errMsg.(string))
		}
	}

	return nil
}
