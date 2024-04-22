package shortlogger

import (
	"go.uber.org/zap"
)

var Sugar zap.SugaredLogger

func Test1() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	Sugar = *logger.Sugar()
	defer logger.Sync()

}
