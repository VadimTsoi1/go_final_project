package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VadimTsoi1/go_final_project/pkg/db"
)

// POST /api/task/done?id=<id>
// Логика (шаг 7):
// - Если repeat пустой - удаляем задачу.
// - Если repeat задан - вычисляем следующую дату и обновляем только поле date.
// Возвращаем пустой JSON {} при успехе, либо {"error":"..."} при ошибке.
func postTaskDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "id required")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}

	// Нет повтора - просто удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{})
		return
	}

	// Есть повтор - считаем nextDate и обновляем поле date
	cur, _ := time.Parse(dateLayout, task.Date)
	next, err := NextDate(cur, task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{})
}
