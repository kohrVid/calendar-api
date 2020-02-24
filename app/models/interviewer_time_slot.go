package models

type InterviewerTimeSlot struct {
	Id            int `json:"id"`
	InterviewerId int `json:"interviewer_id"`
	TimeSlotId    int `json:"time_slot_id"`
}
