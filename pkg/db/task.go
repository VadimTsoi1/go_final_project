package db

import "fmt"

// Task — модель строки таблицы scheduler.
type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// AddTask вставляет задачу в таблицу scheduler и возвращает её id.
func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("db not initialized")
	}
	res, err := DB.Exec(
		`INSERT INTO scheduler(date, title, comment, repeat) VALUES(?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
