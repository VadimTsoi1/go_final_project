package api

import "net/http"

// Регистрируем обработчик /api/task/done при импорте пакета.
func init() {
    http.HandleFunc("/api/task/done", postTaskDone)
}
