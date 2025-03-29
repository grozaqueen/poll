package errs

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type GetErrorCode interface {
	Get(error) (error, int)
}

var (
	InvalidJSONFormat   = errors.New("Неверный формат JSON")
	InvalidReqFormat    = errors.New("Неверный формат запроса")
	InternalServerError = errors.New("Внутренняя ошибка сервера")
	FailedToParseConfig = errors.New("Ошибка парсинга конфигурации")
	MethodNotAllowed    = errors.New("Метод не поддерживается")

	PollNotFound       = errors.New("Опрос не найден")
	PollAlreadyClosed  = errors.New("Голосование уже завершено")
	PollClosedEarly    = errors.New("Голосование завершено досрочно")
	InvalidPollDate    = errors.New("Неверный формат даты опроса")
	PollDateInPast     = errors.New("Дата окончания должна быть в будущем")
	PollCreationFailed = errors.New("Ошибка создания опроса")
	PollDeletionFailed = errors.New("Ошибка удаления опроса")

	UserAlreadyVoted       = errors.New("Вы уже голосовали в этом опросе")
	InvalidVoteOption      = errors.New("Неверный вариант для голосования")
	VoteRegistrationFailed = errors.New("Ошибка регистрации голоса")

	UnauthorizedAccess = errors.New("Недостаточно прав")
	UserNotCreator     = errors.New("Только создатель может выполнить это действие")
	UserNotFound       = errors.New("Пользователь не найден")

	InvalidResponseFromTarantool = errors.New("Некорректный ответ из тарантула")
)

type ErrorStore struct {
	mux        sync.RWMutex
	errorCodes map[error]int
}

func (e *ErrorStore) Get(err error) (error, int) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	errCode, present := e.errorCodes[err]
	if !present {
		log.Println(fmt.Errorf("unexpected error occurred: %w", err))
		return InternalServerError, http.StatusInternalServerError
	}

	return err, errCode
}

func NewErrorStore() *ErrorStore {
	return &ErrorStore{
		mux: sync.RWMutex{},
		errorCodes: map[error]int{
			InvalidJSONFormat:   http.StatusBadRequest,
			InvalidReqFormat:    http.StatusBadRequest,
			InternalServerError: http.StatusInternalServerError,
			FailedToParseConfig: http.StatusInternalServerError,
			MethodNotAllowed:    http.StatusMethodNotAllowed,

			PollNotFound:       http.StatusNotFound,
			PollAlreadyClosed:  http.StatusForbidden,
			PollClosedEarly:    http.StatusForbidden,
			InvalidPollDate:    http.StatusBadRequest,
			PollDateInPast:     http.StatusBadRequest,
			PollCreationFailed: http.StatusInternalServerError,
			PollDeletionFailed: http.StatusInternalServerError,

			UserAlreadyVoted:       http.StatusForbidden,
			InvalidVoteOption:      http.StatusBadRequest,
			VoteRegistrationFailed: http.StatusInternalServerError,

			UnauthorizedAccess: http.StatusUnauthorized,
			UserNotCreator:     http.StatusForbidden,
			UserNotFound:       http.StatusNotFound,

			InvalidResponseFromTarantool: http.StatusInternalServerError,
		},
	}
}
