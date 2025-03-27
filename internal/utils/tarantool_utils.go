package utils

import (
	"errors"
	"github.com/grozaqueen/poll/errs"
	"github.com/tarantool/go-tarantool/v2"
	"log/slog"
)

type TarantoolUtils struct {
	log         *slog.Logger
	errResolver errs.GetErrorCode
}

func NewTarantoolUtils(log *slog.Logger, errResolver errs.GetErrorCode) *TarantoolUtils {
	return &TarantoolUtils{
		log:         log,
		errResolver: errResolver,
	}
}

func (tu *TarantoolUtils) ProcessCall(conn *tarantool.Connection, fnName string, args []interface{}, context string) ([]interface{}, error) {
	resp, err := conn.Call(fnName, args)
	if err != nil {
		tu.log.Error(context+": ошибка вызова Tarantool",
			slog.String("function", fnName),
			slog.String("error", err.Error()))
		return nil, errors.Join(errs.InvalidResponseFromTarantool, err)
	}

	if len(resp) == 0 {
		tu.log.Error(context + ": пустой ответ от Tarantool")
		return nil, errs.InvalidResponseFromTarantool
	}

	return resp, nil
}

func (tu *TarantoolUtils) ExtractMap(resp []interface{}, context string) (map[interface{}]interface{}, error) {
	result, ok := resp[0].(map[interface{}]interface{})
	if !ok {
		tu.log.Error(context + ": некорректный формат ответа Tarantool")
		return nil, errs.InvalidResponseFromTarantool
	}
	return result, nil
}

func (tu *TarantoolUtils) HandleTarantoolError(result map[interface{}]interface{}, context string) error {
	if errMsg, exists := result["error"]; exists {
		errStr, ok := errMsg.(string)
		if !ok {
			tu.log.Error(context + ": неизвестный формат ошибки")
			return errors.New("unknown error format")
		}

		tu.log.Error(context+": ошибка Tarantool",
			slog.String("error", errStr))

		switch errStr {
		case "PollNotFound":
			return errs.PollNotFound
		case "PollAlreadyClosed":
			return errs.PollAlreadyClosed
		case "UserNotCreator":
			return errs.UserNotCreator
		default:
			return errors.New(errStr)
		}
	}
	return nil
}

func (tu *TarantoolUtils) ProcessInsert(conn *tarantool.Connection, space string, data []interface{}, context string) ([]interface{}, error) {
	resp, err := conn.Insert(space, data)
	if err != nil {
		tu.log.Error(context+": ошибка вставки в Tarantool",
			slog.String("space", space),
			slog.String("error", err.Error()))
		return nil, errors.Join(errs.InvalidResponseFromTarantool, err)
	}

	if len(resp) == 0 {
		tu.log.Error(context + ": пустой ответ от Tarantool")
		return nil, errs.InvalidResponseFromTarantool
	}

	return resp, nil
}

func (tu *TarantoolUtils) ExtractTuple(resp []interface{}, context string, expectedFields int) ([]interface{}, error) {
	data, ok := resp[0].([]interface{})
	if !ok || len(data) < expectedFields {
		tu.log.Error(context+": некорректный формат кортежа",
			slog.Int("expected_fields", expectedFields),
			slog.Int("actual_fields", len(data)))
		return nil, errs.InvalidResponseFromTarantool
	}
	return data, nil
}

func (tu *TarantoolUtils) ProcessVoteCall(conn *tarantool.Connection, fnName string, args []interface{}, context string) (map[interface{}]interface{}, error) {
	resp, err := conn.Call(fnName, args)
	if err != nil {
		tu.log.Error(context+": ошибка вызова",
			slog.String("function", fnName),
			slog.String("error", err.Error()))
		return nil, errors.Join(errs.InvalidResponseFromTarantool, err)
	}

	if len(resp) == 0 {
		tu.log.Error(context + ": пустой ответ")
		return nil, errs.InvalidResponseFromTarantool
	}

	result, ok := resp[0].(map[interface{}]interface{})
	if !ok {
		tu.log.Error(context + ": некорректный формат ответа")
		return nil, errs.InvalidResponseFromTarantool
	}

	return result, nil
}

func (tu *TarantoolUtils) HandleVoteError(result map[interface{}]interface{}, context string) error {
	if errMsg, exists := result["error"]; exists {
		errStr, ok := errMsg.(string)
		if !ok {
			tu.log.Error(context + ": неизвестный формат ошибки")
			return errors.New("unknown error format")
		}

		tu.log.Error(context+": ошибка голосования",
			slog.String("error", errStr))

		switch errStr {
		case "PollNotFound":
			return errs.PollNotFound
		case "PollAlreadyClosed":
			return errs.PollAlreadyClosed
		case "UserAlreadyVoted":
			return errs.UserAlreadyVoted
		case "InvalidVoteOption":
			return errs.InvalidVoteOption
		default:
			return errors.New(errStr)
		}
	}
	return nil
}
