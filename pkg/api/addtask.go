package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VadimTsoi1/go_final_project/pkg/db"
)

type addOK struct {
	ID string `json:"id"`
}
type addErr struct {
	Error string `json:"error"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, addErr{Error: err.Error()})
		return
	}

	// title обязателен
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		writeJSON(w, http.StatusBadRequest, addErr{Error: "title required"})
		return
	}

	// "сегодня" — ЛОКАЛЬНАЯ дата, без времени
	now := todayLocal()

	// дата: пустая -> сегодня; задана -> должна быть корректной 20060102
	t.Date = strings.TrimSpace(t.Date)
	if t.Date == "" {
		t.Date = now.Format(dateLayout)
	} else {
		if !isValidDate(t.Date) {
			writeJSON(w, http.StatusBadRequest, addErr{Error: "bad date"})
			return
		}
	}

	// repeat: допускаем "", "y", "d N" (1..400); остальное — ошибка
	t.Repeat = strings.TrimSpace(t.Repeat)
	if t.Repeat != "" && !isValidRepeat(t.Repeat) {
		writeJSON(w, http.StatusBadRequest, addErr{Error: "bad repeat"})
		return
	}

	// сравнение с сегодняшним днём (ЛОКАЛЬНЫЙ)
	tDate, _ := time.ParseInLocation(dateLayout, t.Date, time.Local)
	tDate = normalizeLocal(tDate)

	if tDate.Before(now) { // дата в прошлом
		if t.Repeat == "" {

			t.Date = now.Format(dateLayout)
		} else {

			next, err := NextDate(now, t.Date, t.Repeat)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, addErr{Error: "bad repeat"})
				return
			}
			t.Date = next
		}
	} else {

		if t.Repeat != "" {
			if _, err := NextDate(now, t.Date, t.Repeat); err != nil {
				writeJSON(w, http.StatusBadRequest, addErr{Error: "bad repeat"})
				return
			}
		}
	}

	// вставка в БД
	id, err := db.AddTask(&t)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, addErr{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, addOK{ID: strconv.FormatInt(id, 10)})
}

//helpers

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func isValidDate(s string) bool {
	if len(s) != 8 {
		return false
	}
	_, err := time.ParseInLocation(dateLayout, s, time.Local)
	return err == nil
}

func isValidRepeat(rep string) bool {
	rep = strings.TrimSpace(rep)
	if rep == "" {
		return true
	}
	if rep == "y" {
		return true
	}
	parts := strings.Fields(rep)
	if len(parts) == 2 && parts[0] == "d" {
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return false
		}
		return n >= 1 && n <= 400
	}
	return false // w/m не поддерживаем на базовом шаге
}

// ЛОКАЛЬНЫЙ "сегодня"
func todayLocal() time.Time {
	t := time.Now()
	return normalizeLocal(t)
}

func normalizeLocal(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
