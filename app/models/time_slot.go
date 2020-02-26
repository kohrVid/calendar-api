package models

type TimeSlot struct {
	Id        int    `json:"id"`
	Date      string `json:"date"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
}
