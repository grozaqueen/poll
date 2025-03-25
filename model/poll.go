package model

import "time"

type Poll struct {
	ID       string
	Question string
	Options  []string
	Votes    map[string]int
	EndDate  time.Time
	Creator  struct {
		ID   string
		Name string
	}
	Voters map[string]bool
}
