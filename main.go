package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Poll struct {
	Question string         `json:"question"`
	Options  []string       `json:"options"`
	Votes    map[string]int `json:"votes"`
	EndDate  time.Time      `json:"end_date"`
	Creator  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"creator"`
	Voters map[string]bool `json:"-"`
}

var (
	polls   = make(map[string]*Poll)
	pollsMu sync.RWMutex
)

func main() {
	http.HandleFunc("/poll", createPoll)
	http.HandleFunc("/vote", createVote)
	http.HandleFunc("/results", getResults)
	http.HandleFunc("/poll/complete", completePollEarly)
	http.HandleFunc("/poll/delete", deletePollHandler)
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
