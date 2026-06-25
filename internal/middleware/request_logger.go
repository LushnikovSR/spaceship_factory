package middleware

import (
	"fmt"
	"net/http"
	"time"
)

//RequestLogger создаёт middleware для логирования времени выполнения запроса
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Засекаем время начала обработки запроса
		startTime := time.Now()

		//Логируем начало запроса
		fmt.Printf("⏱️ Начало запроса: %s %s", r.Method, r.URL.Path)

		//Передаём управление следующему обработчику
		next.ServeHTTP(w, r)

		//Вычисляем время выполнения запроса
		duration := time.Since(startTime)

		//Логируем окончание запроса с указанием времени выполнения
		fmt.Printf("✅ Запрос завершен: %s %s, время выполнения: %v", r.Method, r.URL.Path, duration)
	})
}
