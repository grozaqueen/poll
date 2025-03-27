package vote

import (
	"github.com/grozaqueen/poll/internal/model"
	"strconv"
)

func (vr *VoteRepository) CreateVote(vote model.Vote) error {
	const context = "VoteRepository.CreateVote"

	userIDStr := strconv.FormatUint(vote.UserID, 10)
	args := []interface{}{vote.PollID, vote.Option, userIDStr}

	result, err := vr.tarantoolUtils.ProcessVoteCall(vr.polRep.Conn, "create_vote", args, context)
	if err != nil {
		return err
	}

	return vr.tarantoolUtils.HandleVoteError(result, context)
}
