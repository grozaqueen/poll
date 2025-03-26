package model

import "time"

type Poll struct {
	ID       uint64
	Question string
	Options  []string
	Votes    map[string]int
	EndDate  time.Time
	Creator  struct {
		ID   uint64
		Name string
	}
	Voters map[string]bool
}
