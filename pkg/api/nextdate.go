package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

// NextDate вычисляет дату следующего выполнения.
// now     — опорная дата
// dstart  — исходная дата в формате 20060102
// repeat  — правило: "d N" (1..400) или "y"
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	// нормализуем now до полуночи в UTC
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	start, err := time.ParseInLocation(dateLayout, dstart, time.UTC)
	if err != nil {
		return "", fmt.Errorf("bad dstart: %w", err)
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "y":
		// ежегодно + по году, пока не станет > now
		curr := start
		for {
			curr = addYearSpecial(curr)
			if curr.After(now) {
				return curr.Format(dateLayout), nil
			}
		}

	case "d":
		if len(parts) != 2 {
			return "", errors.New("bad repeat format for d")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", errors.New("bad days interval")
		}
		curr := start
		for {
			curr = curr.AddDate(0, 0, n)
			if curr.After(now) {
				return curr.Format(dateLayout), nil
			}
		}

	default:

		return "", errors.New("unsupported repeat format")
	}
}

// addYearSpecial: AddDate(1,0,0) + фиксация 29 фев → 1 мар
func addYearSpecial(t time.Time) time.Time {
	next := t.AddDate(1, 0, 0)
	// Если исходная была 29 фев, то Go даст 28 фев в невисокосный год.
	// По примеру в методичке 1 марта.
	if t.Month() == time.February && t.Day() == 29 && next.Month() == time.February && next.Day() == 28 {
		return next.AddDate(0, 0, 1)
	}
	return next
}

//HTTP

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// GET-параметры: now, date, repeat
	q := r.URL.Query()
	nowStr := strings.TrimSpace(q.Get("now"))
	dstart := strings.TrimSpace(q.Get("date"))
	repeat := strings.TrimSpace(q.Get("repeat"))

	var now time.Time
	if nowStr == "" {
		// текущая дата (UTC) без времени
		t := time.Now().UTC()
		now = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		var err error
		now, err = time.ParseInLocation(dateLayout, nowStr, time.UTC)
		if err != nil {
			http.Error(w, "bad now", http.StatusBadRequest)
			return
		}
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}

	if dstart == "" || repeat == "" {
		http.Error(w, "date/repeat required", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(next))
}
