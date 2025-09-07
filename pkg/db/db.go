package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT,
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	if dbFile == "" {
		return errors.New("empty db file path")
	}

	install := false
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", dbFile)
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}

	if install {
		if _, err := conn.Exec(schema); err != nil {
			_ = conn.Close()
			return fmt.Errorf("apply schema: %w", err)
		}
	}

	DB = conn
	return nil
}
