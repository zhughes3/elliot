package log

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	CallKey      = "call"
	DurationKey  = "dur"
	ErrorKey     = "err"
	SourceKey    = "source"
	MessageKey   = "msg"
	TraceIDKey   = "traceId"
	TimestampKey = "ts"
)

const (
	CalledMessage   = "execution complete"
	TimestampFormat = "2006-01-02T15:04:05.999Z"
)

type zerologger struct {
	logger zerolog.Logger
}

type zerologEntry struct {
	event *zerolog.Event
}

func NewZeroLogger(writer io.Writer) Logger {
	zerolog.TimeFieldFormat = TimestampFormat
	zerolog.TimestampFieldName = TimestampKey
	zerolog.MessageFieldName = MessageKey
	zerolog.DurationFieldInteger = true
	zerolog.CallerFieldName = SourceKey
	zerolog.CallerMarshalFunc = marshalCaller
	zerolog.ErrorFieldName = ErrorKey
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return zerologger{
		logger: zerolog.New(writer).With().CallerWithSkipFrameCount(3).Stack().Timestamp().Logger(),
	}
}

func marshalCaller(_ uintptr, file string, line int) string {
	return afterLastSlash(file) + ":" + strconv.Itoa(line)
}
func afterLastSlash(s string) string {
	return s[strings.LastIndex(s, "/")+1:]
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func sprintlnn(args ...any) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

func (z zerologger) Debug(args ...any) {
	z.logger.Debug().Msg(fmt.Sprint(args...))
}

func (z zerologger) Debugln(args ...any) {
	z.logger.Debug().Msg(sprintlnn(args...))
}

func (z zerologger) Debugf(format string, args ...interface{}) {
	z.logger.Debug().Msgf(format, args...)
}

func (z zerologger) Info(args ...any) {
	z.logger.Info().Msg(fmt.Sprint(args...))
}

func (z zerologger) Infoln(args ...any) {
	z.logger.Info().Msg(sprintlnn(args...))
}

func (z zerologger) Infof(format string, args ...interface{}) {
	z.logger.Info().Msgf(format, args...)
}

func (z zerologger) Warn(args ...any) {
	z.logger.Warn().Msg(fmt.Sprint(args...))
}

func (z zerologger) Warnln(args ...any) {
	z.logger.Warn().Msg(sprintlnn(args...))
}

func (z zerologger) Warnf(format string, args ...interface{}) {
	z.logger.Warn().Msgf(format, args...)
}

func (z zerologger) Error(args ...any) {
	z.logger.Error().Msg(fmt.Sprint(args...))
}

func (z zerologger) Errorln(args ...any) {
	z.logger.Error().Msg(sprintlnn(args...))
}

func (z zerologger) Errorf(format string, args ...interface{}) {
	z.logger.Error().Msgf(format, args...)
}

func (z zerologger) Fatal(args ...any) {
	z.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (z zerologger) Fatalln(args ...any) {
	z.logger.Fatal().Msg(sprintlnn(args...))
}

func (z zerologger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal().Msgf(format, args...)
}
