package commands

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	structs "github.com/fatih/structs"
	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
)

func CreateCandidate(candidate *models.Candidate) (models.Candidate, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	err := db.Insert(candidate)

	if err != nil {
		c := models.Candidate{}

		return c, err
	}

	c := *candidate

	return c, nil
}

func UpdateCandidate(candidate *models.Candidate, params models.Candidate) models.Candidate {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	c := structs.New(candidate)
	p := structs.New(params)

	sql := fmt.Sprintf(
		"UPDATE candidates %v WHERE id = %v;",
		dbHelpers.SetSqlColumns(c, p),
		candidate.Id,
	)

	_, err := db.Model(c).Exec(sql)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	cc := *candidate

	return cc
}

func DeleteCandidate(candidate *models.Candidate) error {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	err := db.Delete(candidate)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return err
}

func CreateCandidateTimeSlot(cid string, timeSlot *models.TimeSlot) (models.TimeSlot, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	ts := models.TimeSlot{}
	cid2, err := strconv.Atoi(cid)

	if err != nil {
		return ts, err
	}

	_, err = db.Query(
		&ts,
		`
	  SELECT 
	      ts.id,
	      ts.start_time,
	      ts.duration
	    FROM time_slots ts
	      INNER JOIN candidate_time_slots cts
	        ON ts.id = cts.time_slot_id
	      INNER JOIN candidates c
	        ON c.id = cts.candidate_id
	      WHERE c.id = ?
	        AND ts.start_time = ?
		AND ts.duration = ?;
	    `,
		cid2,
		timeSlot.StartTime,
		timeSlot.Duration,
	)

	if err != nil {
		log.Println(err)
	}

	if (models.TimeSlot{}) != ts {
		err = errors.New("ERROR #23505 time slot already exists for candidate")

		return ts, err
	}

	ts, err = CreateTimeSlot(timeSlot)

	if err != nil {
		return ts, err
	}

	cts := models.CandidateTimeSlot{
		CandidateId: cid2,
		TimeSlotId:  ts.Id,
	}

	err = db.Insert(&cts)

	if err != nil {
		return ts, err
	}

	return ts, nil
}

func DeleteCandidateTimeSlot(cid string, timeSlot *models.TimeSlot) error {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()
	cid2, err := strconv.Atoi(cid)

	if err != nil {
		return err
	}

	cts := models.CandidateTimeSlot{
		CandidateId: cid2,
		TimeSlotId:  timeSlot.Id,
	}

	err = db.Model(&cts).Select()

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	err = db.Delete(&cts)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	err = DeleteTimeSlot(timeSlot)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return err
}
