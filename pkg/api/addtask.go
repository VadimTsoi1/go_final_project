package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/VadimTsoi1/go_final_project/pkg/db"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

const dateLayout = "20060102"

type addOK struct {
	ID string `json:"id"`
}

type addErr struct {
	Error string `json:"error"`
}

var (
	repeatYear = regexp.MustCompile(`^y$`)
	repeatDay  = regexp.MustCompile(`^d\s+-?\d+$`)
	repeatWeek = regexp.MustCompile(`^w\s+[\d,\s]+$`)
	repeatMon  = regexp.MustCompile(`^m(\s+[\d,\s]+(\s+[\d,\s-]+)?)?$`)
)

func validRepeat(s string) bool {
	if s == "" {
		return true
	}
	s = strings.TrimSpace(s)
	return repeatYear.MatchString(s) ||
		repeatDay.MatchString(s) ||
		repeatWeek.MatchString(s) ||
		repeatMon.MatchString(s)
}

func todayLocal() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ожидаем JSON с полями-строками
	var in struct {
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, addErr{Error: "bad json"})
		return
	}

	title := strings.TrimSpace(in.Title)
	if title == "" {
		writeJSON(w, http.StatusBadRequest, addErr{Error: "title required"})
		return
	}

	// дата: "today" или YYYYMMDD
	var dateStr string
	if strings.EqualFold(strings.TrimSpace(in.Date), "today") {
		dateStr = todayLocal().Format(dateLayout)
	} else {
		if _, err := time.Parse(dateLayout, in.Date); err != nil {
			writeJSON(w, http.StatusBadRequest, addErr{Error: "bad date"})
			return
		}
		dateStr = in.Date
	}

	// repeat — минимальная валидность
	if !validRepeat(in.Repeat) {
		writeJSON(w, http.StatusBadRequest, addErr{Error: "bad repeat"})
		return
	}

	id, err := db.AddTask(&db.Task{
		Date:    dateStr,
		Title:   title,
		Comment: in.Comment,
		Repeat:  strings.TrimSpace(in.Repeat),
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, addOK{ID: id})
}
