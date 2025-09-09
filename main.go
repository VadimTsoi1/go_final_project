package main

import (
    "log"
    "net/http"
    "os"

    "github.com/VadimTsoi1/go_final_project/pkg/api"
    "github.com/VadimTsoi1/go_final_project/pkg/db"
)

func main() {
    // Каталог со статическими файлами веб-интерфейса
    webDir := "./web"

    // Порт HTTP: берём из TODO_PORT, иначе 7540
    port := os.Getenv("TODO_PORT")
    if port == "" {
        port = "7540"
    }
    dbFile := os.Getenv("TODO_DB")
    if dbFile == "" {
        dbFile = "./scheduler.db"
    }

    // Инициализируем БД (при install=true создаём схему)
    if err := db.Init(dbFile, true); err != nil {
        log.Fatalf("db init: %v", err)
    }

    // Регистрируем HTTP-обработчики API
    api.Init()

    // Отдаём статические файлы из ./web
    fs := http.FileServer(http.Dir(webDir))
    http.Handle("/", fs)

    addr := ":" + port
    log.Printf("server listening on %s, db=%s", addr, dbFile)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

