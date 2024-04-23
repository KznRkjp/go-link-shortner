package shortlogger

import (
	"fmt"
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
		fmt.Println("here")

		// отправляем сведения о запросе в zap
		// defer Sugar.Sync()
		fmt.Println(uri, method, duration)

		sugar.Infow(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
		// defer Sugar.Sync()

	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
