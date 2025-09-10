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
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    TEXT NOT NULL,
    title   TEXT NOT NULL,
    comment TEXT NOT NULL DEFAULT '',
    repeat  TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// Init открывает файл БД SQLite и при install=true применяет схему.
func Init(filename string, install bool) error {
    if filename == "" {
        return errors.New("empty db filename")
    }

    if err := os.MkdirAll(".", 0o755); err != nil {
        return fmt.Errorf("mkdir: %w", err)
    }

    conn, err := sql.Open("sqlite", filename)
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

// Close closes the global DB connection.
func Close() error {
    if DB != nil {
        err := DB.Close()
        DB = nil
        return err
    }
    return nil
}
