package models

type CandidateTimeSlot struct {
	Id          int `json:"id"`
	CandidateId int `json:"candidate_id"`
	TimeSlotId  int `json:"time_slot_id"`
}
