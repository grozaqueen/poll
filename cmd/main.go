package main

import (
	"github.com/gorilla/mux"
	pollDeliveryLib "github.com/grozaqueen/poll/delivery/poll"
	voteDeliveryLib "github.com/grozaqueen/poll/delivery/vote"
	errResolveLib "github.com/grozaqueen/poll/errs"
	pollRepoLib "github.com/grozaqueen/poll/repository/poll"
	voteRepoLIb "github.com/grozaqueen/poll/repository/vote"
	pollUsecaseLib "github.com/grozaqueen/poll/usecase/poll"
	"log"
	"net/http"
)

type pollDelivery interface {
	CompletePollEarly(w http.ResponseWriter, r *http.Request)
	CreatePoll(w http.ResponseWriter, r *http.Request)
	DeletePoll(w http.ResponseWriter, r *http.Request)
	GetResults(w http.ResponseWriter, r *http.Request)
}

type voteDelivery interface {
	CreateVote(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	router *mux.Router
	polls  pollDelivery
	votes  voteDelivery
}

func NewServer() (*Server, error) {
	router := mux.NewRouter()
	errResolver := errResolveLib.NewErrorStore()

	pollRepo := pollRepoLib.NewPollRepository()
	pollService := pollUsecaseLib.NewPollUseCase(pollRepo)
	pollDeliv := pollDeliveryLib.NewPollDelivery(pollService, pollRepo, errResolver)

	voteRepo := voteRepoLIb.NewVoteRepository(pollRepo)
	voteDeliv := voteDeliveryLib.NewVoteDelivery(voteRepo, errResolver)

	return &Server{
		router: router,
		polls:  pollDeliv,
		votes:  voteDeliv,
	}, nil
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/poll", server.polls.CreatePoll)
	http.HandleFunc("/vote", server.votes.CreateVote)
	http.HandleFunc("/results", server.polls.GetResults)
	http.HandleFunc("/poll/complete", server.polls.CompletePollEarly)
	http.HandleFunc("/poll/delete", server.polls.DeletePoll)
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
