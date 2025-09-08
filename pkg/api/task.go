package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/VadimTsoi1/go_final_project/pkg/db"
)

// /api/task: GET, POST, PUT
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// добавление — как на шаге 4
		addTaskHandler(w, r)

	case http.MethodGet:
		// получить одну по id
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			writeJSONError(w, http.StatusBadRequest, "id required")
			return
		}
		t, err := db.GetTask(idStr)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "not found")
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id":      strconv.FormatInt(t.ID, 10),
			"date":    t.Date,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		})

	case http.MethodPut:
		// обновление (валидируем так же, как при POST)
		var in struct {
			ID      string `json:"id"`
			Date    string `json:"date"`
			Title   string `json:"title"`
			Comment string `json:"comment"`
			Repeat  string `json:"repeat"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSONError(w, http.StatusBadRequest, "bad json")
			return
		}
		if in.ID == "" {
			writeJSONError(w, http.StatusBadRequest, "id required")
			return
		}
		// проверка id — число
		id, err := strconv.ParseInt(in.ID, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "bad id")
			return
		}
		// проверка title
		if in.Title == "" {
			writeJSONError(w, http.StatusBadRequest, "title required")
			return
		}
		// проверка date
		if _, err := timeParseDate(in.Date); err != nil {
			writeJSONError(w, http.StatusBadRequest, "bad date")
			return
		}
		// проверка repeat
		if !validRepeat(in.Repeat) {
			writeJSONError(w, http.StatusBadRequest, "bad repeat")
			return
		}

		if err := db.UpdateTask(&db.Task{
			ID:      id,
			Date:    in.Date,
			Title:   in.Title,
			Comment: in.Comment,
			Repeat:  in.Repeat,
		}); err != nil {
			// если id не найден
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		// успешное изменение → пустой JSON {}
		_ = json.NewEncoder(w).Encode(map[string]any{})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

// timeParseDate — обёртка, чтобы не тянуть сюда todayLocal()
func timeParseDate(s string) (string, error) {
	_, err := time.Parse(dateLayout, s)
	if err != nil {
		return "", err
	}
	return s, nil
}
