package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"go.opentelemetry.io/otel/trace"

	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
)

type contextKey string

const (
	skipLogKey = contextKey("skip-log")
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger(ctx context.Context) *Logger {
	cfg := env.GetConfig()
	var l zerolog.Logger
	if cfg.AppEnv == "local" {
		// Use ConsoleWriter for human-friendly, colorized output in local dev.
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		l = zerolog.New(output).With().Timestamp().Logger()
	} else {
		l = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	logger := &Logger{
		logger: l,
	}
	logger = logger.setTrace(ctx)
	return logger
}

func InitializeLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// MarshalStackが errs.Err の StackTrace() メソッドを呼び出してくれます。
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func (l *Logger) setTrace(ctx context.Context) *Logger {
	cfg := env.GetConfig()
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	if !sc.IsValid() {
		return l
	}

	traceStr := fmt.Sprintf("projects/%s/traces/%s", cfg.GCPProjectID, sc.TraceID().String())
	l.logger = l.logger.With().
		Str("logging.googleapis.com/trace", traceStr).
		Str("logging.googleapis.com/spanId", sc.SpanID().String()).Logger()

	return l
}

// func (l *Logger) setUserInfo(ctx context.Context) *Logger {
// 	if v, ok := ctx.Value("CPID").(string); ok {
// 		l.logger = l.logger.
// 			With().Str("CPID", v).Logger()
// 	}

// 	if v, ok := ctx.Value("userID").(string); ok {
// 		l.logger = l.logger.
// 			With().Str("userID", v).Logger()
// 	}

// 	return l
// }

func (l *Logger) DebugWriter() zerolog.Logger {
	return l.logger.With().Str("severity", "DEBUG").Logger()
}

func (l *Logger) DebugLog(msg string) {
	l.logger.Debug().
		Str("severity", "DEBUG").
		Msg(msg)
}

func (l *Logger) InfoEvent() *zerolog.Event {
	return l.logger.Info().Str("severity", "INFO")
}

func (l *Logger) InfoLog(msg string) {
	l.InfoEvent().Msg(msg)
}

func (l *Logger) SkipLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipLogKey, true)
}

func (l *Logger) ShouldSkipLog(ctx context.Context) bool {
	return ctx.Value(skipLogKey) != nil
}

func (l *Logger) WarnLog(err error) {
	var e *errs.Err
	event := l.logger.Warn().Stack()
	if errors.As(err, &e) {
		event.Str("result_code", strconv.Itoa(int(e.ResultCode())))
	}
	event.Err(err).Str("severity", "WARNING").Msg("")
}

func (l *Logger) ErrorLog(err error) {
	var e *errs.Err
	event := l.logger.Error().Stack()
	if errors.As(err, &e) {
		event.Str("result_code", strconv.Itoa(int(e.ResultCode())))
	}
	event.Err(err).Str("severity", "ERROR").Msg("")
}

func (l *Logger) FatalLog(err error) {
	var e *errs.Err
	event := l.logger.Fatal().Stack()
	if errors.As(err, &e) {
		event.Str("result_code", strconv.Itoa(int(e.ResultCode())))
	}
	// Cloud LoggingにはFatalレベルのログがないため、ALERTレベルのログを出力する
	event.Err(err).Str("severity", "ALERT").Msg("")
}
