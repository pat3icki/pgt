package server

import (
	"log/slog"
)

// Logger defines the interface for logging operations.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

type logger struct {
	log *slog.Logger
}

func DefaultLogger() Logger {
	return logger{
		log: slog.Default(),
	}
}

func (l logger) Debug(msg string, args ...any) {}
func (l logger) Info(msg string, args ...any)  {}
func (l logger) Warn(msg string, args ...any)  {}
func (l logger) Error(msg string, args ...any) {}
func (l logger) Fatal(msg string, args ...any) {}
