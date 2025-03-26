package poll

import (
	"errors"
	"fmt"
	"github.com/grozaqueen/poll/internal/model"
	"github.com/grozaqueen/poll/internal/utils"
	"time"
)

func (pr *PollRepository) GetResults(pollID uint64) (model.Poll, error) {
	// Вызываем функцию в Tarantool
	resp, err := pr.Call("get_poll_results", []interface{}{pollID})
	if err != nil {
		return model.Poll{}, fmt.Errorf("tarantool call failed: %v", err)
	}

	// Проверяем ответ
	if len(resp) == 0 {
		return model.Poll{}, errors.New("poll not found")
	}

	// Преобразуем ответ
	result, ok := resp[0].(map[interface{}]interface{})
	if !ok {
		return model.Poll{}, fmt.Errorf("invalid response format: expected map, got %T", resp[0])
	}

	// Обрабатываем голоса
	rawVotes, votesOk := result["votes"]
	votes := make(map[string]int)
	if votesOk {
		if votesMap, ok := rawVotes.(map[interface{}]interface{}); ok {
			votes = utils.ConvertToIntMap(votesMap)
		}
	}

	// Обрабатываем дату окончания
	var endDate time.Time
	if rawEndDate, exists := result["end_date"]; exists {
		if endDateUnix, err := toInt64(rawEndDate); err == nil {
			endDate = time.Unix(endDateUnix, 0)
		}
	}

	// Обрабатываем информацию о создателе
	var creator struct {
		ID   uint64
		Name string
	}
	if rawCreatorID, exists := result["creator_id"]; exists {
		creator.ID = utils.InterfaceToUint64(rawCreatorID)
	}
	if rawCreatorName, exists := result["creator_name"]; exists {
		creator.Name = toString(rawCreatorName)
	}

	// Обрабатываем список проголосовавших
	rawVoters, _ := result["voters"]
	voters := utils.ConvertToBoolMap(rawVoters)

	// Собираем результат
	poll := model.Poll{
		ID:       utils.InterfaceToUint64(result["id"]),
		Question: toString(result["question"]),
		Options:  utils.ConvertToStringSlice(result["options"]),
		Votes:    votes,
		EndDate:  endDate,
		Creator:  creator,
		Voters:   voters,
	}

	return poll, nil
}

// Вспомогательные функции
func toString(val interface{}) string {
	if s, ok := val.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", val)
}
