package poll

import (
	"github.com/grozaqueen/poll/model"
	"sync"
)

type PollRepository struct {
	Polls   map[string]*model.Poll
	PollsMu sync.RWMutex
}

func NewPollRepository() *PollRepository {
	return &PollRepository{
		Polls:   make(map[string]*model.Poll),
		PollsMu: sync.RWMutex{},
	}
}
