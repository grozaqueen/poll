package main

import (
	"github.com/gorilla/mux"
	errResolveLib "github.com/grozaqueen/poll/errs"
	pollDeliveryLib "github.com/grozaqueen/poll/internal/delivery/poll"
	voteDeliveryLib "github.com/grozaqueen/poll/internal/delivery/vote"
	"github.com/grozaqueen/poll/internal/logger"
	pollRepoLib "github.com/grozaqueen/poll/internal/repository/poll"
	voteRepoLIb "github.com/grozaqueen/poll/internal/repository/vote"
	pollUsecaseLib "github.com/grozaqueen/poll/internal/usecase/poll"
	"github.com/grozaqueen/poll/internal/utils"
	"github.com/grozaqueen/poll/middleware"
	"github.com/tarantool/go-tarantool/v2"

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
	router        *mux.Router
	polls         pollDelivery
	votes         voteDelivery
	tarantoolConn *tarantool.Connection
}

func (s *Server) SetupRoutes() {
	s.router.HandleFunc("/poll", s.polls.CreatePoll).Methods("POST")
	s.router.HandleFunc("/vote", s.votes.CreateVote).Methods("POST")
	s.router.HandleFunc("/results", s.polls.GetResults).Methods("GET")
	s.router.HandleFunc("/poll/complete", s.polls.CompletePollEarly).Methods("POST")
	s.router.HandleFunc("/poll/delete", s.polls.DeletePoll).Methods("DELETE")
}

func NewServer() (*Server, error) {
	tarantoolConn, err := InitTarantool()
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	errResolver := errResolveLib.NewErrorStore()

	log := logger.InitLogger()

	router.Use(middleware.RequestLogger(log))

	utilsHelpers := utils.NewHandlerUtils(log, errResolver)
	tarantoolUtils := utils.NewTarantoolUtils(log, errResolver)

	pollRepo := pollRepoLib.NewPollRepository(tarantoolConn, tarantoolUtils, log)
	pollService := pollUsecaseLib.NewPollUseCase(pollRepo, log, errResolver)
	pollDeliv := pollDeliveryLib.NewPollDelivery(pollService, pollRepo, utilsHelpers)

	voteRepo := voteRepoLIb.NewVoteRepository(pollRepo, tarantoolUtils)
	voteDeliv := voteDeliveryLib.NewVoteDelivery(voteRepo, utilsHelpers)

	return &Server{
		router:        router,
		polls:         pollDeliv,
		votes:         voteDeliv,
		tarantoolConn: tarantoolConn,
	}, nil
}

func (s *Server) Close() {
	if s.tarantoolConn != nil {
		s.tarantoolConn.Close()
	}
}
