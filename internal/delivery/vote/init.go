package vote

import (
	"github.com/grozaqueen/poll/errs"
	"github.com/grozaqueen/poll/internal/model"
)

type VoteCreator interface {
	CreateVote(vote model.Vote) error
}

type VoteDelivery struct {
	VoteCreator VoteCreator
	errResolver errs.GetErrorCode
}

func NewVoteDelivery(VoteCreator VoteCreator,
	errResolver errs.GetErrorCode) *VoteDelivery {
	return &VoteDelivery{
		VoteCreator: VoteCreator,
		errResolver: errResolver,
	}
}
