package poll

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/model"
	"github.com/grozaqueen/poll/internal/utils"
	"log/slog"
	"time"
)

func (pm *PollUseCase) CreatePoll(poll model.Poll) (model.Poll, error) {
	nowInMSK := time.Now().In(utils.MskLocation)
	if poll.EndDate.Before(nowInMSK) {
		err, _ := pm.errResolver.Get(errs.PollDateInPast)
		pm.log.Error("[ PollUsecase. CreatePoll ] Указанное время голосования должно быть в будущем", slog.String("error", err.Error()))
		return model.Poll{}, errs.PollDateInPast
	}
	poll, err := pm.PollRepository.CreatePoll(poll)
	if err != nil {
		pm.log.Error("[ PollUsecase. CreatePoll ] Ошибка на уровне репозитория", slog.String("error", err.Error()))
		return model.Poll{}, err
	}
	return poll, err
}
