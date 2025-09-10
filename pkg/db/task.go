package db

import (
	"database/sql"
	"fmt"
)

// Task - модель задачи, строка в таблице scheduler.
type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// AddTask добавляет задачу и возвращает её id (как строку).
func AddTask(t *Task) (string, error) {
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

// Tasks возвращает будущие задачи (date >= fromDate),
// отсортированные по date и id, с ограничением limit.
func Tasks(fromDate string, limit int) ([]*Task, error) {
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
    // возвращаем nil, если записей нет (вызывающая сторона формирует JSON-модель)
	return out, nil
}

// TasksSearch ищет задачи по подстроке в title/comment.
func TasksSearch(search string, limit int) ([]*Task, error) {
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
    // возвращаем nil при отсутствии результатов
	return out, nil
}

// TasksOnDate возвращает задачи на конкретную дату YYYYMMDD.
func TasksOnDate(dateYYYYMMDD string, limit int) ([]*Task, error) {
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
    // возвращаем nil при отсутствии результатов
	return out, nil
}

// GetTask возвращает задачу по id (ошибка, если не найдена).
func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	var t Task
	if err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTask обновляет поля задачи по ID.
func UpdateTask(task *Task) error {
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

// TruncateAll полностью очищает таблицу (для тестов).
func TruncateAll(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM scheduler`)
	return err
}

// DeleteTask удаляет задачу по id.
func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("incorrect id for deleting task")
	}
	return nil
}

// UpdateDate обновляет только поле date у задачи.
func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date=? WHERE id=?`, next, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("incorrect id for updating date")
	}
	return nil
}
