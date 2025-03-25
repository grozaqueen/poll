package vote

import "github.com/grozaqueen/poll/model"

type CreateVoteRequest struct {
	PollID string `json:"poll_id"`
	Option string `json:"option"`
	UserID string `json:"user_id"`
}

func (r CreateVoteRequest) toModel() model.Vote {
	return model.Vote{
		PollID: r.PollID,
		Option: r.Option,
		UserID: r.UserID,
	}
}
