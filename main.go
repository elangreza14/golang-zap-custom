package main

import (
	"errors"

	"github.com/elangreza14/golang-zap-custom/logger"
)

func main() {
	v := errors.New("cek")
	// y := http.ErrAbortHandler
	// opt := logger.Option{}
	opt := &logger.Option{
		EncodingType: logger.EncodingTypeConsole,
		NameService:  "CEK",
	}
	a, _ := logger.NewLogger(opt)

	a.Info("cek 0", "qqwqwq", nil, opt)
	a.Error("cek 1", v)
	a.Debug("cek 2", v)
}
