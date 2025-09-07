package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VadimTsoi1/go_final_project/pkg/api"
	"github.com/VadimTsoi1/go_final_project/pkg/db"
)

func main() {
	//web
	webDir := "./web"
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" //порт по умолчанию
	}

	//db
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	//регистрация API
	api.Init()

	// файловый сервер, отдает из ./web
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	//запуск
	addr := ":" + port
	log.Printf("server listening on %s, db=%s", addr, dbFile)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
