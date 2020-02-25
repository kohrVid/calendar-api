package models

type TimeSlot struct {
	Id        int    `json:"id"`
	Date      string `json:"date"`
	StartTime int    `json:"start_time"`
	Duration  int    `json:"duration"`
}
