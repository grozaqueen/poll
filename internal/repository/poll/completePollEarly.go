package poll

import (
	"errors"
	"fmt"
	"github.com/grozaqueen/poll/errs"
	"time"
)

func (pr *PollRepository) CompletePollEarly(pollID uint64, userID uint64) (time.Time, error) {
	resp, err := pr.Call("complete_poll_early", []interface{}{pollID, userID})
	if err != nil {
		return time.Time{}, fmt.Errorf("tarantool call failed: %v", err)
	}

	if len(resp) == 0 {
		return time.Time{}, errors.New("empty response from Tarantool")
	}

	result, ok := resp[0].(map[interface{}]interface{})
	if !ok {
		return time.Time{}, errors.New("invalid response format")
	}

	// Обработка ошибок
	if errMsg, exists := result["error"]; exists {
		if errStr, ok := errMsg.(string); ok {
			switch errStr {
			case "PollNotFound":
				return time.Time{}, errs.PollNotFound
			case "UserNotCreator":
				return time.Time{}, errs.UserNotCreator
			default:
				return time.Time{}, errors.New(errStr)
			}
		}
		return time.Time{}, errors.New("unknown error format")
	}

	// Получаем новую дату окончания
	endDateUnix, err := toInt64(result["end_date"])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid end_date format: %v", err)
	}

	return time.Unix(endDateUnix, 0), nil
}
