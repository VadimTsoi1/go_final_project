package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/VadimTsoi1/go_final_project/pkg/db"
)

type taskOut struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type tasksResp struct {
	Tasks []taskOut `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const yyyymmdd = "20060102"
	limit := 50
	search := r.URL.Query().Get("search")

	var (
		items []*db.Task
		err   error
	)

	if search != "" {
		if t, e := time.Parse("02.01.2006", search); e == nil {
			items, err = db.TasksOnDate(t.Format(yyyymmdd), limit)
		} else {
			items, err = db.TasksSearch(search, limit)
		}
	} else {
		today := time.Now().Format(yyyymmdd)
		items, err = db.Tasks(today, limit)
	}

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	out := tasksResp{Tasks: make([]taskOut, 0, len(items))}
	for _, t := range items {
		out.Tasks = append(out.Tasks, taskOut{
			ID:      strconv.FormatInt(t.ID, 10),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}
