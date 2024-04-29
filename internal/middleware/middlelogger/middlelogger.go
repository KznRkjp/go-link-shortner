package middlelogger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func ServerStartLog(addr string) {
	// создаём предустановленный регистратор zap
	var sugar1 zap.SugaredLogger
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	sugar1 = *logger.Sugar()
	defer logger.Sync()
	sugar1.Infow(
		"Starting server",
		"addr", addr,
	)

}

func ServerErrorLog(error string) {
	// создаём предустановленный регистратор zap
	var sugar1 zap.SugaredLogger
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	sugar1 = *logger.Sugar()
	defer logger.Sync()
	sugar1.Error(
		"Server error",
		"error", error,
	)

}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		logger, err := zap.NewDevelopment()
		if err != nil {
			// вызываем панику, если ошибка
			panic(err)
		}
		sugar := *logger.Sugar()
		// функция Now() возвращает текущее время
		start := time.Now()

		// эндпоинт /ping
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		// точка, где выполняется хендлер pingHandler

		h.ServeHTTP(w, r) // обслуживание оригинального запроса

		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
		duration := time.Since(start)

		// отправляем сведения о запросе в zap
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)

	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
