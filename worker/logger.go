package worker

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Debug(args ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Info(args ...interface{}) {
	log.WithLevel(zerolog.InfoLevel).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Warn(args ...interface{}) {
	log.WithLevel(zerolog.WarnLevel).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Error(args ...interface{}) {
	log.WithLevel(zerolog.ErrorLevel).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Fatal(args ...interface{}) {
	log.WithLevel(zerolog.FatalLevel).Msg(fmt.Sprint(args...))
}
