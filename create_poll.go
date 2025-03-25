package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var mskLocation = time.FixedZone("MSK", 3*60*60)

func createPoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	type CreatePollRequest struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
		EndDate  string   `json:"end_date"`
		UserID   string   `json:"user_id"`
		UserName string   `json:"username"`
	}

	type CreatePollResponse struct {
		PollID   string   `json:"poll_id"`
		Options  []string `json:"options"`
		UserID   string   `json:"user_id"`
		UserName string   `json:"username"`
	}

	var req CreatePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	endDate, err := parseSimpleDate(req.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nowInMSK := time.Now().In(mskLocation)
	if endDate.Before(nowInMSK) {
		http.Error(w, fmt.Sprintf("Дата окончания должна быть в будущем. Текущее время: %s",
			nowInMSK.Format("02.01.2006 15:04:05")),
			http.StatusBadRequest)
		return
	}

	pollID := generateID()

	pollsMu.Lock()
	polls[pollID] = &Poll{
		Question: req.Question,
		Options:  req.Options,
		Votes:    make(map[string]int),
		EndDate:  endDate,
		Creator: struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			ID:   req.UserID,
			Name: req.UserName,
		},
	}
	pollsMu.Unlock()

	response := CreatePollResponse{
		PollID:   pollID,
		Options:  req.Options,
		UserID:   req.UserID,
		UserName: req.UserName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func generateID() string {
	return "poll_" + strconv.Itoa(len(polls)+1)
}

func parseSimpleDate(dateStr string) (time.Time, error) {
	if len(dateStr) != 10 || dateStr[2] != '.' || dateStr[5] != '.' {
		return time.Time{}, fmt.Errorf("неверный формат даты, используйте DD.MM.YYYY")
	}

	day, err1 := strconv.Atoi(dateStr[:2])
	month, err2 := strconv.Atoi(dateStr[3:5])
	year, err3 := strconv.Atoi(dateStr[6:])

	if err1 != nil || err2 != nil || err3 != nil {
		return time.Time{}, fmt.Errorf("дата должна содержать только цифры и точки")
	}

	if day < 1 || day > 31 || month < 1 || month > 12 || year < 2000 {
		return time.Time{}, fmt.Errorf("некорректная дата")
	}

	return time.Date(
		year,
		time.Month(month),
		day,
		23, 59, 59, 0,
		mskLocation,
	), nil
}
