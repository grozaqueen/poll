package vote

import (
	"github.com/grozaqueen/poll/internal/model"
)

type CreateVoteRequest struct {
	PollID uint64 `json:"poll_id"`
	Option string `json:"option"`
	UserID uint64 `json:"user_id"`
}

func (r CreateVoteRequest) toModel() model.Vote {
	return model.Vote{
		PollID: r.PollID,
		Option: r.Option,
		UserID: r.UserID,
	}
}
