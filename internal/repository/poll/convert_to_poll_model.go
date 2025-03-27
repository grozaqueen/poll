package poll

import (
	"github.com/grozaqueen/poll/internal/model"
	"github.com/grozaqueen/poll/internal/utils"
	"time"
)

func (pr *PollRepository) convertToPollModel(data map[interface{}]interface{}) (model.Poll, error) {
	votes := make(map[string]int)
	if rawVotes, exists := data["votes"]; exists {
		if votesMap, ok := rawVotes.(map[interface{}]interface{}); ok {
			votes = utils.ConvertToIntMap(votesMap)
		}
	}

	var endDate time.Time
	if rawEndDate, exists := data["end_date"]; exists {
		if endDateUnix, err := utils.ToInt64(rawEndDate); err == nil {
			endDate = time.Unix(endDateUnix, 0)
		}
	}

	var creator model.Creator
	if rawCreatorID, exists := data["creator_id"]; exists {
		creator.ID = utils.InterfaceToUint64(rawCreatorID)
	}
	if rawCreatorName, exists := data["creator_name"]; exists {
		creator.Name = utils.ToString(rawCreatorName)
	}

	voters := make(map[string]bool)
	if rawVoters, exists := data["voters"]; exists {
		voters = utils.ConvertToBoolMap(rawVoters)
	}

	return model.Poll{
		ID:       utils.InterfaceToUint64(data["id"]),
		Question: utils.ToString(data["question"]),
		Options:  utils.ConvertToStringSlice(data["options"]),
		Votes:    votes,
		EndDate:  endDate,
		Creator:  creator,
		Voters:   voters,
	}, nil
}
