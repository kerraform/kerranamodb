package logging

import (
	"fmt"
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level string

const (
	LevelInfo        Level = "info"
	LevelDebug       Level = "debug"
	LevelWarn        Level = "warn"
	LevelError       Level = "error"
	LevelUnspecified Level = ""
)

type Format string

const (
	FormatConsole      Format = "console"
	FormatColorConsole Format = "color"
	FormatJSON         Format = "json"
	FormatUnspecified  Format = ""
)

var (
	consoleEncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "M",
		LevelKey:       "L",
		TimeKey:        "T",
		NameKey:        "N",
		CallerKey:      "C",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	colorConsoleEncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "M",
		LevelKey:       "L",
		TimeKey:        "T",
		NameKey:        "N",
		CallerKey:      "C",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	jsonEncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
)

func NewLogger(writer io.Writer, level Level, format Format) (*zap.Logger, error) {
	zapLevel, err := convertZapLevel(level)
	if err != nil {
		return nil, err
	}
	zapEncoder, err := convertZapEncoder(format)
	if err != nil {
		return nil, err
	}
	return zap.New(
		zapcore.NewCore(
			zapEncoder,
			zapcore.Lock(zapcore.AddSync(writer)),
			zap.NewAtomicLevelAt(zapLevel),
		),
	), nil
}

func convertZapLevel(level Level) (zapcore.Level, error) {
	l := strings.TrimSpace(strings.ToLower(string(level)))
	switch Level(l) {
	case LevelInfo, LevelUnspecified:
		return zapcore.InfoLevel, nil
	case LevelDebug:
		return zapcore.DebugLevel, nil
	case LevelWarn:
		return zapcore.WarnLevel, nil
	case LevelError:
		return zapcore.ErrorLevel, nil
	default:
		return 0, fmt.Errorf("unknown log level [info,debug,warn,error]: %q", level)
	}
}

func convertZapEncoder(format Format) (zapcore.Encoder, error) {
	f := strings.TrimSpace(strings.ToLower(string(format)))
	switch Format(f) {
	case FormatColorConsole, FormatUnspecified:
		return zapcore.NewConsoleEncoder(colorConsoleEncoderConfig), nil
	case FormatConsole:
		return zapcore.NewConsoleEncoder(consoleEncoderConfig), nil
	case FormatJSON:
		return zapcore.NewJSONEncoder(jsonEncoderConfig), nil
	default:
		return nil, fmt.Errorf("unknown log format [console,color,json]: %q", format)
	}
}
