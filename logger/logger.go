package logger

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger interface {
		Info(logStr string, payload ...interface{})
		Error(logStr string, payload ...interface{})
		Debug(logStr string, payload ...interface{})
	}

	logger struct {
		zapBase *zap.Logger
	}
)

func NewLogger(opt *Option) (Logger, error) {

	config := zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		DisableCaller:    true,
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "time",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},

		Level: zap.NewAtomicLevelAt(zapcore.DebugLevel),
	}

	if opt != nil && opt.EncodingType == EncodingTypeConsole {
		config.Encoding = string(EncodingTypeConsole)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if opt != nil && opt.EnableStackTrace {
		config.EncoderConfig.StacktraceKey = "stacktrace"
	}

	if opt != nil && opt.NameService != "" {
		appName := make(map[string]interface{})
		appName["service"] = opt.NameService
		config.InitialFields = appName
	}

	zapBuilder, err := config.Build()

	if err != nil {
		return nil, err
	}

	return &logger{
		zapBase: zapBuilder,
	}, nil
}

func setField(params ...interface{}) []zap.Field {
	res := []zap.Field{}

	for i := 0; i < len(params); i++ {
		typePrams := reflect.TypeOf(params[i])

		if typePrams == nil {
			index := fmt.Sprintf("nil-%v", i+1)
			res = append(res, zap.Any(index, params[i]))
			continue
		}

		switch typePrams.Kind().String() {
		case "chan":
			index := fmt.Sprintf("chan-%v", i+1)
			res = append(res, zap.Any(index, "cannot log chan"))
		case "ptr":
			switch typePrams.String() {
			case "*errors.errorString":
				index := fmt.Sprintf("error-%v", i+1)
				res = append(res, zap.Any(index, params[i]))
			default:
				index := fmt.Sprintf("pointer-%v", i+1)
				res = append(res, zap.Any(index, params[i]))
			}
		default:
			index := fmt.Sprintf("data-%v", i+1)
			res = append(res, zap.Any(index, params[i]))
		}
	}

	return res
}

func (l *logger) Info(logStr string, payload ...interface{}) {
	base := setField(payload...)
	l.zapBase.With(base...).Info(logStr)
}

func (l *logger) Error(logStr string, payload ...interface{}) {
	base := setField(payload...)
	l.zapBase.With(base...).Error(logStr)
}

func (l *logger) Debug(logStr string, payload ...interface{}) {
	base := setField(payload...)
	l.zapBase.With(base...).Debug(logStr)
}
