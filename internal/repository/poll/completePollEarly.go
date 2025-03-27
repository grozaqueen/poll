package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/utils"
	"log/slog"
	"time"
)

func (pr *PollRepository) CompletePollEarly(pollID uint64, userID uint64) (time.Time, error) {
	const context = "CompletePollEarly"

	resp, err := pr.tarantoolUtils.ProcessCall(pr.Conn, "complete_poll_early", []interface{}{pollID, userID}, context)
	if err != nil {
		return time.Time{}, err
	}

	result, err := pr.tarantoolUtils.ExtractMap(resp, context)
	if err != nil {
		return time.Time{}, err
	}

	if err := pr.tarantoolUtils.HandleTarantoolError(result, context); err != nil {
		return time.Time{}, err
	}

	endDateUnix, err := utils.ToInt64(result["end_date"])
	if err != nil {
		pr.log.Error(context+": некорректная дата окончания",
			slog.String("error", err.Error()))
		return time.Time{}, errs.InvalidResponseFromTarantool
	}

	return time.Unix(endDateUnix, 0), nil
}
