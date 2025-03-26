package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	errResolveLib "github.com/grozaqueen/poll/errs"
	pollDeliveryLib "github.com/grozaqueen/poll/internal/delivery/poll"
	voteDeliveryLib "github.com/grozaqueen/poll/internal/delivery/vote"
	pollRepoLib "github.com/grozaqueen/poll/internal/repository/poll"
	voteRepoLIb "github.com/grozaqueen/poll/internal/repository/vote"
	pollUsecaseLib "github.com/grozaqueen/poll/internal/usecase/poll"
	"github.com/tarantool/go-tarantool/v2"
	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/decimal"
	_ "github.com/tarantool/go-tarantool/v2/uuid"
	"log"
	"net/http"
	"time"
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

func initTarantool() (*tarantool.Connection, error) {
	// Увеличиваем время ожидания инициализации Tarantool
	time.Sleep(5 * time.Second)

	// Увеличиваем таймауты подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  "tarantool_db:3301", // Жестко задаем адрес вместо переменной окружения
		User:     "adminka",           // Жестко задаем пользователя
		Password: "12345",             // Жестко задаем пароль
	}

	opts := tarantool.Opts{
		Timeout:       5 * time.Second, // Увеличили таймаут
		Reconnect:     1 * time.Second,
		MaxReconnects: 5, // Увеличили количество попыток
	}

	// Подключаемся с новыми параметрами
	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}

	// Проверяем не только пинг, но и существование пространства polls
	resp, err := conn.Do(tarantool.NewEvalRequest("return box.space.polls ~= nil")).Get()
	if err != nil {
		return nil, fmt.Errorf("space check failed: %v", err)
	}

	if len(resp) == 0 || !resp[0].(bool) {
		return nil, errors.New("space 'polls' does not exist")
	}

	log.Println("Successfully connected to Tarantool and verified 'polls' space exists")
	return conn, nil
}

func NewServer() (*Server, error) {
	tarantoolConn, err := initTarantool()
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	errResolver := errResolveLib.NewErrorStore()

	pollRepo := pollRepoLib.NewPollRepository(tarantoolConn)
	pollService := pollUsecaseLib.NewPollUseCase(pollRepo)
	pollDeliv := pollDeliveryLib.NewPollDelivery(pollService, pollRepo, errResolver)

	voteRepo := voteRepoLIb.NewVoteRepository(pollRepo)
	voteDeliv := voteDeliveryLib.NewVoteDelivery(voteRepo, errResolver)

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

func (s *Server) SetupRoutes() {
	s.router.HandleFunc("/poll", s.polls.CreatePoll).Methods("POST")
	s.router.HandleFunc("/vote", s.votes.CreateVote).Methods("POST")
	s.router.HandleFunc("/results", s.polls.GetResults).Methods("GET")
	s.router.HandleFunc("/poll/complete", s.polls.CompletePollEarly).Methods("POST")
	s.router.HandleFunc("/poll/delete", s.polls.DeletePoll).Methods("DELETE")
}

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
