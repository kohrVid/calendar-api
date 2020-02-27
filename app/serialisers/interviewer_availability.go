package serialisers

import "github.com/kohrVid/calendar-api/app/models"

type InterviewerAvailability struct {
	InterviewerId int               `json:"interviewer_id"`
	TimeSlots     []models.TimeSlot `json:"time_slots"`
}
