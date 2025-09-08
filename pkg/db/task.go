package db

import (
	"database/sql"
	"fmt"
)

// Task — модель строки таблицы scheduler.
type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// AddTask вставляет новую задачу и возвращает id строкой.
func AddTask(t *Task) (string, error) {
	if DB == nil {
		return "", fmt.Errorf("db not initialized")
	}
	res, err := DB.Exec(
		`INSERT INTO scheduler(date, title, comment, repeat) VALUES(?,?,?,?)`,
		t.Date, t.Title, t.Comment, t.Repeat,
	)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return fmt.Sprint(id), nil
}

// Tasks возвращает ближайшие задачи (дата >= fromDate), сортировка по дате↑, затем id↑, ограничение limit.
func Tasks(fromDate string, limit int) ([]*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if limit <= 0 {
		limit = 50
	}
	rows, err := DB.Query(`
		SELECT id, date, title, comment, repeat
		  FROM scheduler
		 WHERE date >= ?
		 ORDER BY date, id
		 LIMIT ?`, fromDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = make([]*Task, 0)
	}
	return out, nil
}

// TasksSearch возвращает задачи по подстроке в title/comment.
func TasksSearch(search string, limit int) ([]*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if limit <= 0 {
		limit = 50
	}
	like := "%" + search + "%"
	rows, err := DB.Query(`
		SELECT id, date, title, comment, repeat
		  FROM scheduler
		 WHERE title LIKE ? OR comment LIKE ?
		 ORDER BY date, id
		 LIMIT ?`, like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = make([]*Task, 0)
	}
	return out, nil
}

// TasksOnDate возвращает задачи на конкретную дату YYYYMMDD.
func TasksOnDate(dateYYYYMMDD string, limit int) ([]*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if limit <= 0 {
		limit = 50
	}
	rows, err := DB.Query(`
		SELECT id, date, title, comment, repeat
		  FROM scheduler
		 WHERE date = ?
		 ORDER BY date, id
		 LIMIT ?`, dateYYYYMMDD, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = make([]*Task, 0)
	}
	return out, nil
}

// GetTask возвращает задачу по id (строкой из запроса).
func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	var t Task
	if err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTask обновляет все поля по ID.
func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("db not initialized")
	}
	res, err := DB.Exec(
		`UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

// TruncateAll — утилита для тестов
func TruncateAll(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM scheduler`)
	return err
}
