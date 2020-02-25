package queries

import (
	"fmt"
	"strconv"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
)

func ListCandidates() []models.Candidate {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)

	candidates := make([]models.Candidate, 0)
	err := db.Model(&candidates).Select()

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return candidates
}

func FindCandidate(id string) (models.Candidate, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	idx, err := strconv.Atoi(id)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	candidate := &models.Candidate{Id: idx}
	err = db.Select(candidate)

	if err != nil {
		return *candidate, err
	}

	return *candidate, nil
}

func ListCandidateTimeSlots(cid string) []models.TimeSlot {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	timeSlots := make([]models.TimeSlot, 0)

	_, err := db.Query(
		&timeSlots,
		`
	  SELECT
	      ts.id,
	      ts.date,
	      ts.start_time,
	      ts.duration
	    FROM time_slots ts
	    INNER JOIN candidate_time_slots cts
	    ON ts.id = cts.time_slot_id
	    INNER JOIN candidates c
	    ON c.id = cts.candidate_id
	    WHERE c.id = ?`,
		cid,
	)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return timeSlots
}

func FindCandidateTimeSlot(cid string, id string) (models.TimeSlot, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	timeSlot := &models.TimeSlot{}

	_, err := db.QueryOne(
		timeSlot,
		`
	  SELECT
	      ts.id,
	      ts.date,
	      ts.start_time,
	      ts.duration
	    FROM time_slots ts
	    INNER JOIN candidate_time_slots cts
	    ON ts.id = cts.time_slot_id
	    INNER JOIN candidates c
	    ON c.id = cts.candidate_id
	    WHERE c.id = ? AND ts.id = ?`,
		cid,
		id,
	)

	if err != nil {
		return *timeSlot, err
	}

	return *timeSlot, nil
}
