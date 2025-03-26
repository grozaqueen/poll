package poll

import (
	"errors"
	"fmt"
	"github.com/grozaqueen/poll/internal/model"
	"github.com/grozaqueen/poll/internal/utils"
	"math"
	"strconv"
	"time"
)

func (pr *PollRepository) CreatePoll(poll model.Poll) (model.Poll, error) {
	endDateUnix := uint64(poll.EndDate.Unix())
	createdAtUnix := uint64(time.Now().Unix())
	resp, err := pr.Insert("polls", []interface{}{
		nil,
		poll.Question,
		poll.Options,
		make(map[string]int),
		endDateUnix,
		poll.Creator.ID,
		poll.Creator.Name,
		make(map[string]bool),
		createdAtUnix,
	})

	if err != nil {
		return model.Poll{}, err
	}

	if len(resp) == 0 {
		return model.Poll{}, errors.New("empty response from Tarantool")
	}

	data, ok := resp[0].([]interface{})
	if !ok || len(data) < 8 {
		return model.Poll{}, errors.New("invalid response format")
	}

	dataEnd, err := toInt64(data[4])
	if err != nil {
		return model.Poll{}, errors.New("invalid response format")
	}
	createdPoll := model.Poll{
		ID:       utils.InterfaceToUint64(data[0]),
		Question: data[1].(string),
		Options:  utils.ConvertToStringSlice(data[2]),
		Votes:    utils.ConvertToIntMap(data[3]),
		EndDate:  time.Unix(dataEnd, 0),
		Creator: struct {
			ID   uint64
			Name string
		}{
			ID:   utils.InterfaceToUint64(data[5]),
			Name: data[6].(string),
		},
		Voters: utils.ConvertToBoolMap(data[7]),
	}

	return createdPoll, nil
}

func toInt64(val interface{}) (int64, error) {
	if val == nil {
		return 0, errors.New("nil value")
	}

	switch v := val.(type) {
	case int64:
		return v, nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, errors.New("value too large for int64")
		}
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case float64:
		if v > math.MaxInt64 || v < math.MinInt64 {
			return 0, errors.New("value out of int64 range")
		}
		return int64(v), nil
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string as int64: %v", err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported type for int64: %T", val)
	}
}
