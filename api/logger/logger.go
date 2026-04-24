// Package logger expone un logger estructurado global construido sobre log/slog.
// Usa salida JSON en producción; texto legible en todos los demás entornos.
//
// Uso:
//
//	logger.L().Info("servidor iniciado", "port", port)
//	logger.L().With("request_id", rid).Error("db falló", "err", err)
package logger

import (
	"log/slog"
	"os"
)

var l *slog.Logger

func init() {
	setup()
}

func setup() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}

	var handler slog.Handler
	if os.Getenv("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	l = slog.New(handler)
	slog.SetDefault(l)
}

// L retorna el logger global.
func L() *slog.Logger {
	return l
}
