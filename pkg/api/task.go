package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "strings"

    "github.com/VadimTsoi1/go_final_project/pkg/db"
)

// Маршрут /api/task: GET, POST, PUT, DELETE
func taskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        // Создание задачи (шаг 4)
        addTaskHandler(w, r)

    case http.MethodGet:
        // Получение задачи по id
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
        // Обновление задачи (формат такой же, как у POST)
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
        // Проверяем и парсим id
        id, err := strconv.ParseInt(in.ID, 10, 64)
        if err != nil {
            writeJSONError(w, http.StatusBadRequest, "bad id")
            return
        }
        // Проверяем title
        if in.Title == "" {
            writeJSONError(w, http.StatusBadRequest, "title required")
            return
        }
        // Проверяем date
        if _, err := timeParseDate(in.Date); err != nil {
            writeJSONError(w, http.StatusBadRequest, "bad date")
            return
        }
        // Проверяем repeat (поддерживаем только базовые правила)
        if !validRepeat(in.Repeat) {
            writeJSONError(w, http.StatusBadRequest, "bad repeat")
            return
        }
        // Если дата в прошлом и repeat задан — переносим на следующую дату
        if strings.TrimSpace(in.Repeat) != "" {
            now := todayLocal().Format(dateLayout)
            if in.Date < now {
                next, err := NextDate(todayLocal(), in.Date, strings.TrimSpace(in.Repeat))
                if err != nil {
                    writeJSONError(w, http.StatusBadRequest, "bad repeat")
                    return
                }
                in.Date = next
            }
        }

        if err := db.UpdateTask(&db.Task{
            ID:      id,
            Date:    in.Date,
            Title:   in.Title,
            Comment: in.Comment,
            Repeat:  in.Repeat,
        }); err != nil {
            // Некорректный id при обновлении
            writeJSONError(w, http.StatusBadRequest, err.Error())
            return
        }
        // Успех: возвращаем пустой JSON {}
        _ = json.NewEncoder(w).Encode(map[string]any{})

    case http.MethodDelete:
        idStr := r.URL.Query().Get("id")
        if idStr == "" {
            writeJSONError(w, http.StatusBadRequest, "id required")
            return
        }
        if err := db.DeleteTask(idStr); err != nil {
            writeJSONError(w, http.StatusBadRequest, err.Error())
            return
        }
        _ = json.NewEncoder(w).Encode(map[string]any{})

    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }

}

// Парсим дату в формате YYYYMMDD (локаль todayLocal не трогаем)
func timeParseDate(s string) (string, error) {
    _, err := time.Parse(dateLayout, s)
    if err != nil {
        return "", err
    }
    return s, nil
}
