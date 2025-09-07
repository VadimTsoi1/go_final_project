package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VadimTsoi1/go_final_project/pkg/api"
	"github.com/VadimTsoi1/go_final_project/pkg/db"
)

func main() {
	// статические файлы
	webDir := "./web"

	// порт и БД из окружения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	dbFile := os.Getenv("TODO_DB")
	if dbFile == "" {
		dbFile = "./scheduler.db"
	}

	// init DB (+ схема)
	if err := db.Init(dbFile, true); err != nil {
		log.Fatalf("db init: %v", err)
	}

	// регистрируем API
	api.Init()

	// файловый сервер из ./web
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	addr := ":" + port
	log.Printf("server listening on %s, db=%s", addr, dbFile)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
