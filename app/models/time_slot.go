package models

type TimeSlot struct {
	Id        int `json:"id"`
	StartTime int `json:"start_time"`
	Duration  int `json:"duration"`
}