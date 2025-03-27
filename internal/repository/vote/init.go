package vote

import (
	"github.com/grozaqueen/poll/internal/repository/poll"
	"github.com/tarantool/go-tarantool/v2"
)

type tarantoolUtils interface {
	ProcessVoteCall(conn *tarantool.Connection, fnName string, args []interface{}, context string) (map[interface{}]interface{}, error)
	HandleVoteError(result map[interface{}]interface{}, context string) error
}

type VoteRepository struct {
	polRep         *poll.PollRepository
	tarantoolUtils tarantoolUtils
}

func NewVoteRepository(pollrepository *poll.PollRepository, tarantoolUtils tarantoolUtils) *VoteRepository {
	return &VoteRepository{
		polRep:         pollrepository,
		tarantoolUtils: tarantoolUtils,
	}
}
