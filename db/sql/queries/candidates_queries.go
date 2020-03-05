package queries

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/app/serialisers"
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

func ListCandidateTimeSlots(cid string, all ...bool) []models.TimeSlot {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	timeSlots := make([]models.TimeSlot, 0)
	currentTime := "NOW()"

	if (os.Getenv("ENV") == "test") && (os.Getenv("TIME_NOW") != "") {
		currentTime = os.Getenv("TIME_NOW")
	}

	sql := `
	  SELECT
	      ts.id,
	      ts.date,
	      ts.start_time,
	      ts.end_time
	    FROM time_slots ts
	    INNER JOIN candidate_time_slots cts
	      ON ts.id = cts.time_slot_id
	    INNER JOIN candidates c
	      ON c.id = cts.candidate_id
	    WHERE c.id = ?`

	if !(all != nil && all[0]) {
		sql += fmt.Sprintf(`
	AND CONCAT(ts.date, ' ', ts.start_time, ':00:00+00')::timestamp >= %v;`,
			currentTime,
		)
	}

	_, err := db.Query(
		&timeSlots,
		sql,
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
	      ts.end_time
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

func ListCandidateAndInterviewerTimeSlot(
	cid string,
	interviewers []string,
	all ...bool,
) []serialisers.InterviewerAvailability {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	currentTime := "NOW()"
	a := false

	if (os.Getenv("ENV") == "test") && (os.Getenv("TIME_NOW") != "") {
		currentTime = os.Getenv("TIME_NOW")
	}

	if all != nil && all[0] {
		a = true
	}

	availability := make([]serialisers.InterviewerAvailability, 0)
	interviewersIds := []int{}

	sql := `
	  SELECT
	      ts.id,
	      ts.date,
	      ts.start_time,
	      CASE WHEN ts.end_time > ? THEN ?
		ELSE ts.end_time
	      END end_time
	    FROM time_slots ts
	    INNER JOIN interviewer_time_slots its
	      ON ts.id = its.time_slot_id
	    INNER JOIN interviewers i
	      ON i.id = its.interviewer_id
	    WHERE i.id = ?`

	if !a {
		sql += fmt.Sprintf(`
	AND CONCAT(ts.date, ' ', ts.start_time, ':00:00+00')::timestamp >= %v`,
			currentTime,
		)
	}

	sql += `
	AND (CONCAT(ts.date, ' ', ts.start_time, ':00:00+00')::timestamp >= CONCAT(?, ' ', ?, ':00:00+00')::timestamp
	  AND CONCAT(ts.date, ' ', ts.start_time, ':00:00+00')::timestamp <= CONCAT(?, ' ', ?, ':00:00+00')::timestamp);`

	for idx, iid := range interviewers {
		id, err := strconv.Atoi(iid)

		if err != nil {
			fmt.Errorf("Error: %v", err)
		}

		interviewersIds = append(interviewersIds, id)
		candidateAvailability := ListCandidateTimeSlots(cid, a)
		timeSlots := make([]models.TimeSlot, 0)

		if len(candidateAvailability) > 0 {
			for _, ca := range candidateAvailability {
				ts := make([]models.TimeSlot, 0)

				_, err := db.Query(
					&ts,
					sql,
					ca.EndTime,
					ca.EndTime,
					interviewers[idx],
					ca.Date,
					ca.StartTime,
					ca.Date,
					ca.EndTime,
				)

				if err != nil {
					fmt.Errorf("Error: %v", err)
				}

				timeSlots = append(timeSlots, ts...)
			}

		}

		availability = append(
			availability,
			serialisers.InterviewerAvailability{
				InterviewerId: interviewersIds[idx],
				TimeSlots:     timeSlots,
			},
		)
	}

	return availability
}
