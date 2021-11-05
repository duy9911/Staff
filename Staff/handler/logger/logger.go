package logger

import (
	"errors"
	"reflect"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func Logger(title string, val interface{}) {
	sugar := logger.Sugar()
	new := errors.New("sas")

	typeVal := reflect.TypeOf(val).Kind()
	typeErr := reflect.TypeOf(new).Kind()

	if typeVal == typeErr {
		sugar.Errorf(title, val)
		return
	}
	sugar.Infof(title, val)
}
