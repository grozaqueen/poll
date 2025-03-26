package poll

import (
	"github.com/tarantool/go-tarantool/v2"
)

type PollRepository struct {
	Conn *tarantool.Connection
}

func NewPollRepository(conn *tarantool.Connection) *PollRepository {
	return &PollRepository{
		Conn: conn,
	}
}

func (pr *PollRepository) Insert(space string, tuple []interface{}) ([]interface{}, error) {
	req := tarantool.NewInsertRequest(space).Tuple(tuple)
	return pr.Conn.Do(req).Get()
}

func (pr *PollRepository) Call(functionName string, args []interface{}) ([]interface{}, error) {
	req := tarantool.NewCallRequest(functionName).Args(args)
	return pr.Conn.Do(req).Get()
}
