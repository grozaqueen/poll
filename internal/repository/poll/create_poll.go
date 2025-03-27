package poll

import (
	"errors"
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/model"
	"github.com/grozaqueen/poll/internal/utils"
	"log/slog"
	"time"
)

func (pr *PollRepository) CreatePoll(poll model.Poll) (model.Poll, error) {
	const context = "PollRepository.CreatePoll"

	endDateUnix := uint64(poll.EndDate.Unix())
	createdAtUnix := uint64(time.Now().Unix())

	insertData := []interface{}{
		nil,
		poll.Question,
		poll.Options,
		make(map[string]int),
		endDateUnix,
		poll.Creator.ID,
		poll.Creator.Name,
		make(map[string]bool),
		createdAtUnix,
	}

	resp, err := pr.Conn.Insert("polls", insertData)
	if err != nil {
		pr.log.Error(context+": ошибка вставки",
			slog.String("error", err.Error()))
		return model.Poll{}, errors.Join(errs.InvalidResponseFromTarantool, err)
	}

	if len(resp) == 0 {
		pr.log.Error(context + ": пустой ответ")
		return model.Poll{}, errs.InvalidResponseFromTarantool
	}

	data, ok := resp[0].([]interface{})
	if !ok || len(data) < 8 {
		pr.log.Error(context+": некорректный формат ответа",
			slog.Int("полей получено", len(data)))
		return model.Poll{}, errs.InvalidResponseFromTarantool
	}

	dataEnd, err := utils.ToInt64(data[4])
	if err != nil {
		pr.log.Error(context+": ошибка преобразования даты",
			slog.String("error", err.Error()))
		return model.Poll{}, errs.InvalidResponseFromTarantool
	}

	createdPoll := model.Poll{
		ID:       utils.InterfaceToUint64(data[0]),
		Question: data[1].(string),
		Options:  utils.ConvertToStringSlice(data[2]),
		Votes:    utils.ConvertToIntMap(data[3]),
		EndDate:  time.Unix(dataEnd, 0),
		Creator: model.Creator{
			ID:   utils.InterfaceToUint64(data[5]),
			Name: data[6].(string),
		},
		Voters: utils.ConvertToBoolMap(data[7]),
	}

	return createdPoll, nil
}
