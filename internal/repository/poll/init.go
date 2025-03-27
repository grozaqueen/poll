package poll

import (
	"github.com/tarantool/go-tarantool/v2"
	"log/slog"
)

type tarantoolUtils interface {
	ProcessCall(conn *tarantool.Connection, fnName string, args []interface{}, context string) ([]interface{}, error)
	ExtractMap(resp []interface{}, context string) (map[interface{}]interface{}, error)
	HandleTarantoolError(result map[interface{}]interface{}, context string) error
}
type PollRepository struct {
	Conn           *tarantool.Connection
	tarantoolUtils tarantoolUtils
	log            *slog.Logger
}

func NewPollRepository(conn *tarantool.Connection, tarantoolUtils tarantoolUtils, log *slog.Logger) *PollRepository {
	return &PollRepository{
		Conn:           conn,
		tarantoolUtils: tarantoolUtils,
		log:            log,
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
