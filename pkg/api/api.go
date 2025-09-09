package api

import "net/http"

// Init регистрирует HTTP-обработчики API.
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)   // POST - создание; GET/PUT/DELETE - работа с задачей
	http.HandleFunc("/api/tasks", tasksHandler) // GET - список и поиск
}
