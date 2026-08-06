package logger

import (
	"log/slog"
	"os"
)

//Создание структурированного логгера,
//для валидного покрытия кода логами

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
