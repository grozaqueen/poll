package utils

import (
	"encoding/json"
	"github.com/grozaqueen/poll/errs"
	"log/slog"
	"net/http"
	"strconv"
)

type HandlerUtils struct {
	log         *slog.Logger
	errResolver errs.GetErrorCode
}

func NewHandlerUtils(log *slog.Logger, errResolver errs.GetErrorCode) *HandlerUtils {
	return &HandlerUtils{
		log:         log,
		errResolver: errResolver,
	}
}

func (h *HandlerUtils) ValidateMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		h.HandleError(w, r, errs.MethodNotAllowed, "ValidateMethod: неподдерживаемый метод - "+r.Method)
		return false
	}
	return true
}

func (h *HandlerUtils) DecodeRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		h.HandleError(w, r, errs.InvalidJSONFormat, "Ошибка декодирования запроса")
		return false
	}
	return true
}

func (h *HandlerUtils) SendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		if h.log != nil {
			h.log.Error("Ошибка кодирования ответа",
				slog.String("error", err.Error()),
			)
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *HandlerUtils) HandleError(w http.ResponseWriter, r *http.Request, err error, context string) {
	resolvedErr, status := h.errResolver.Get(err)

	if h.log != nil {
		h.log.Error(context,
			slog.String("error", resolvedErr.Error()),
			slog.Int("status", status),
			slog.String("method", r.Method),
		)
	}

	http.Error(w, resolvedErr.Error(), status)
}

func (h *HandlerUtils) ParseUintParam(w http.ResponseWriter, r *http.Request, param string) (uint64, bool) {
	if param == "" {
		h.HandleError(w, r, errs.InvalidReqFormat, "Отсутствует обязательный параметр")
		return 0, false
	}

	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		h.HandleError(w, r, err, "Неверный формат параметра")
		return 0, false
	}

	return val, true
}
