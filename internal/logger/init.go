package logger

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"
)

type PrettyJSONHandler struct {
	slog.Handler
	output *os.File
}

func (h *PrettyJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	data := map[string]interface{}{
		"time":    r.Time.Format(time.RFC3339Nano),
		"level":   r.Level.String(),
		"message": r.Message,
	}

	r.Attrs(func(attr slog.Attr) bool {
		data[attr.Key] = attr.Value.Any()
		return true
	})

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = h.output.Write(append(jsonData, '\n'))
	return err
}

func InitLogger() *slog.Logger {
	env := os.Getenv("APP_ENV")
	level := slog.LevelDebug

	switch env {
	case "prod":
		level = slog.LevelError
	case "staging":
		level = slog.LevelInfo
	}

	handler := &PrettyJSONHandler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
		output: os.Stdout,
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
