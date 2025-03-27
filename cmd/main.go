package main

import (
	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/decimal"
	_ "github.com/tarantool/go-tarantool/v2/uuid"

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
