package poll

import (
	"github.com/grozaqueen/poll/model"
	"time"
)

type CompletePollRequest struct {
	PollID string `json:"poll_id"`
	UserID string `json:"user_id"`
}

type CompletePollResponse struct {
	Status  string `json:"status"`
	EndDate string `json:"end_date"`
	PollID  string `json:"poll_id"`
}

type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	EndDate  string   `json:"end_date"`
	UserID   string   `json:"user_id"`
	UserName string   `json:"username"`
}

func (r CreatePollRequest) toModel(endDate time.Time) model.Poll {
	return model.Poll{
		Question: r.Question,
		Options:  r.Options,
		EndDate:  endDate,
		Creator: struct {
			ID   string
			Name string
		}{ID: r.UserID, Name: r.UserName},
	}
}

type CreatePollResponse struct {
	PollID   string   `json:"poll_id"`
	Options  []string `json:"options"`
	UserID   string   `json:"user_id"`
	UserName string   `json:"username"`
}

type DeletePollRequest struct {
	PollID string `json:"poll_id"`
	UserID string `json:"user_id"`
}

type DeletePollResponse struct {
	Status    string `json:"status"`
	PollID    string `json:"poll_id"`
	DeletedBy string `json:"deleted_by"`
}

type ResultsPollResponse struct {
	Question string         `json:"question"`
	Options  []string       `json:"options"`
	Votes    map[string]int `json:"votes"`
}
