package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	webDir := "./web"

	// читаем порт из переменной окружения TODO_PORT
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" //порт по умолчанию
	}

	// файловый сервер, отдает из ./web
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
