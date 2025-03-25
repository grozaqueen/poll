package poll

import (
	"github.com/grozaqueen/poll/model"
	"strconv"
)

func (pr *PollRepository) CreatePoll(poll model.Poll) (model.Poll, error) {
	pollID := pr.generateID()

	pr.PollsMu.Lock()
	pollInMap := model.Poll{
		ID:       pollID,
		Question: poll.Question,
		Options:  poll.Options,
		Votes:    make(map[string]int),
		EndDate:  poll.EndDate,
		Creator: struct {
			ID   string
			Name string
		}{
			ID:   poll.Creator.ID,
			Name: poll.Creator.Name,
		},
	}
	pr.Polls[pollID] = &pollInMap
	pr.PollsMu.Unlock()

	return pollInMap, nil
}

func (pr *PollRepository) generateID() string {
	return "poll_" + strconv.Itoa(len(pr.Polls)+1)
}
