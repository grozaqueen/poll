package vote

import (
	"github.com/grozaqueen/poll/repository/poll"
)

type VoteRepository struct {
	polRep *poll.PollRepository
}

func NewVoteRepository(pollrepository *poll.PollRepository) *VoteRepository {
	return &VoteRepository{
		polRep: pollrepository,
	}
}
