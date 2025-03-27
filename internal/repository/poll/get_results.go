package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/model"
	"log/slog"
)

func (pr *PollRepository) GetResults(pollID uint64) (model.Poll, error) {
	const context = "PollRepository.GetResults"

	resp, err := pr.tarantoolUtils.ProcessCall(pr.Conn, "get_poll_results", []interface{}{pollID}, context)
	if err != nil {
		return model.Poll{}, err
	}

	result, err := pr.tarantoolUtils.ExtractMap(resp, context)
	if err != nil {
		return model.Poll{}, errs.PollNotFound
	}

	if err := pr.tarantoolUtils.HandleTarantoolError(result, context); err != nil {
		return model.Poll{}, err
	}

	poll, err := pr.convertToPollModel(result)
	if err != nil {
		pr.log.Error(context+": ошибка преобразования данных",
			slog.String("error", err.Error()))
		return model.Poll{}, errs.InvalidResponseFromTarantool
	}

	return poll, nil
}
