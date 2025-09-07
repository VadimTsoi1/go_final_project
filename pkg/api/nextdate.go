package api

import (
	"errors"
	"net/http"
	"time"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if dstart == "" || repeat == "" {
		return "", errors.New("empty args")
	}

	_, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", err
	}
	return dstart, nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	dstart := q.Get("date")
	repeat := q.Get("repeat")
	if dstart == "" || repeat == "" {
		http.Error(w, "date/repeat required", http.StatusBadRequest)
		return
	}

	now := todayLocal()
	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(next))
}
