package api

import "net/http"

// Init регистрирует все API-эндпоинты
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
}
