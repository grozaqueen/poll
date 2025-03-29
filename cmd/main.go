// @title Poll API
// @version 1.0
// @description API системы опросов и голосований
// @host localhost:8080
// @BasePath /
package main

import (
	_ "github.com/grozaqueen/poll/docs"
	"log"
	"net/http"
)

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	server.SetupRoutes()

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", server.router))
}
