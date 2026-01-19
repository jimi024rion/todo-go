package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
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
	zerolog.ErrorStackMarshaler = errs.MarshalStack
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
	if errors.As(err, &e) {
		l.logger.Warn().
			Str("severity", "WARNING").
			Str("result_code", strconv.Itoa(int(e.ResultCode()))).
			Stack().
			Err(e).
			Msg("")
	} else {
		l.logger.Warn().
			Str("severity", "WARNING").
			Stack().
			Err(err).
			Msg("")
	}
}

func (l *Logger) ErrorLog(err error) {
	var e *errs.Err
	if errors.As(err, &e) {
		l.logger.Error().
			Str("severity", "ERROR").
			Str("result_code", strconv.Itoa(int(e.ResultCode()))).
			Stack().
			Err(e).
			Msg("")
	} else {
		l.logger.Error().
			Str("severity", "ERROR").
			Stack().
			Err(err).
			Msg("")
	}
}

func (l *Logger) FatalLog(err error) {
	var e *errs.Err
	if errors.As(err, &e) {
		// Cloud LoggingにはFatalレベルのログがないため、ALERTレベルのログを出力する
		l.logger.Fatal().
			Str("severity", "ALERT").
			Str("result_code", strconv.Itoa(int(e.ResultCode()))).
			Stack().
			Err(e).
			Msg("")
	} else {
		l.logger.Fatal().
			Str("severity", "ALERT").
			Stack().
			Err(err).
			Msg("")
	}
}
