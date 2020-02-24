package queries

import (
	"fmt"
	"strconv"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
)

func ListInterviewers() []models.Interviewer {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)

	interviewers := make([]models.Interviewer, 0)
	err := db.Model(&interviewers).Select()

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return interviewers
}

func FindInterviewer(id string) (models.Interviewer, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	idx, err := strconv.Atoi(id)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	interviewer := &models.Interviewer{Id: idx}
	err = db.Select(interviewer)

	if err != nil {
		return *interviewer, err
	}

	return *interviewer, nil
}

func ListInterviewerTimeSlots(iid string) []models.TimeSlot {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	timeSlots := make([]models.TimeSlot, 0)

	_, err := db.Query(
		&timeSlots,
		`
	  SELECT
	      ts.id,
	      ts.start_time,
	      ts.duration
	    FROM time_slots ts
	    INNER JOIN interviewer_time_slots its
	    ON ts.id = its.time_slot_id
	    INNER JOIN interviewers i
	    ON i.id = its.interviewer_id
	    WHERE i.id = ?`,
		iid,
	)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return timeSlots
}

func FindInterviewerTimeSlot(iid string, id string) (models.TimeSlot, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	timeSlot := &models.TimeSlot{}

	_, err := db.QueryOne(
		timeSlot,
		`
	  SELECT
	      ts.id,
	      ts.start_time,
	      ts.duration
	    FROM time_slots ts
	    INNER JOIN interviewer_time_slots its
	    ON ts.id = its.time_slot_id
	    INNER JOIN interviewers i
	    ON i.id = its.interviewer_id
	    WHERE i.id = ? AND ts.id = ?`,
		iid,
		id,
	)

	if err != nil {
		return *timeSlot, err
	}

	return *timeSlot, nil
}
