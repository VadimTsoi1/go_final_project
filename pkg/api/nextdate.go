package api

import (
    "errors"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
    if dstart == "" || repeat == "" {
        return "", errors.New("empty args")
    }
    start, err := time.Parse(dateLayout, dstart)
    if err != nil {
        return "", err
    }
    r := strings.TrimSpace(repeat)
    // Повтор каждый год
    if r == "y" {
        base := now
        if now.Before(start) {
            base = start
        }
        y, m, d := start.Date()
        // Кандидат в том же месяце/дне для базового года
        cand := time.Date(base.Year(), m, d, 0, 0, 0, 0, base.Location())
        if !cand.After(base) {
            y = base.Year() + 1
        } else {
            y = base.Year()
        }
        // 29 февраля: переносим на 1 марта, если год не високосный
        cand = time.Date(y, m, d, 0, 0, 0, 0, base.Location())
        if cand.Month() != m || cand.Day() != d {
            cand = time.Date(y, time.March, 1, 0, 0, 0, 0, base.Location())
        }
        return cand.Format(dateLayout), nil
    }
    // Ежедневный повтор: d N
    if strings.HasPrefix(r, "d") {
        parts := strings.Fields(r)
        if len(parts) != 2 {
            return "", errors.New("bad repeat")
        }
        n, err := strconv.Atoi(parts[1])
        if err != nil || n <= 0 || n > 400 {
            return "", errors.New("bad repeat")
        }
        base := now
        if now.Before(start) {
            // Первая дата после старта
            return start.AddDate(0, 0, n).Format(dateLayout), nil
        }
        // Считаем смещение в днях
        diff := int(base.Sub(start).Hours() / 24)
        steps := diff/n + 1
        next := start.AddDate(0, 0, steps*n)
        return next.Format(dateLayout), nil
    }
    return "", errors.New("bad repeat")
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
    if v := r.URL.Query().Get("now"); v != "" {
        if t, err := time.Parse(dateLayout, v); err == nil {
            now = t
        }
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

