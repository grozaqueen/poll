package vote

import (
	"github.com/grozaqueen/poll/internal/repository/poll"
)

type VoteRepository struct {
	polRep *poll.PollRepository
}

func NewVoteRepository(pollrepository *poll.PollRepository) *VoteRepository {
	return &VoteRepository{
		polRep: pollrepository,
	}
}
