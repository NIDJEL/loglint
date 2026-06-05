package loglintcases

import "log/slog"

func bad() {
	token := "abc"

	slog.Info("Starting server") // want "log message should start with lowercase letter"
	slog.Info("failed to connect")
	slog.Info("warning: something went wrong...") // want "log message should not contain special characters or emoji"
	slog.Info("server started 🔥")                 // want "log message should not contain special characters or emoji"
	slog.Info("запуск сервера")                   // want "log message should be written in English"
	slog.Info("token: " + token)                  // want "log message should not contain sensitive data"
}
