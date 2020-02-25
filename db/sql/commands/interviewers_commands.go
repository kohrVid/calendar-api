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

func CreateInterviewer(interviewer *models.Interviewer) (models.Interviewer, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	err := db.Insert(interviewer)

	if err != nil {
		c := models.Interviewer{}

		return c, err
	}

	c := *interviewer

	return c, nil
}

func UpdateInterviewer(interviewer *models.Interviewer, params models.Interviewer) models.Interviewer {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	c := structs.New(interviewer)
	p := structs.New(params)

	sql := fmt.Sprintf(
		"UPDATE interviewers %v WHERE id = %v;",
		dbHelpers.SetSqlColumns(c, p),
		interviewer.Id,
	)

	_, err := db.Model(c).Exec(sql)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	cc := *interviewer

	return cc
}

func DeleteInterviewer(interviewer *models.Interviewer) error {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	err := db.Delete(interviewer)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return err
}

func CreateInterviewerTimeSlot(iid string, timeSlot *models.TimeSlot) (models.TimeSlot, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	ts := models.TimeSlot{}
	iid2, err := strconv.Atoi(iid)

	if err != nil {
		return ts, err
	}

	_, err = db.Query(
		&ts,
		`
	  SELECT 
	      ts.id,
	      ts.date,
	      ts.start_time,
	      ts.end_time
	    FROM time_slots ts
	      INNER JOIN interviewer_time_slots its
	        ON ts.id = its.time_slot_id
	      INNER JOIN interviewers i
	        ON i.id = its.interviewer_id
	      WHERE i.id = ?
	        AND ts.date = ?
	        AND ts.start_time = ?
		AND ts.end_time = ?;
	    `,
		iid2,
		timeSlot.Date,
		timeSlot.StartTime,
		timeSlot.EndTime,
	)

	if err != nil {
		log.Println(err)
	}

	if (models.TimeSlot{}) != ts {
		err = errors.New("ERROR #23505 time slot already exists for interviewer")

		return ts, err
	}

	ts, err = CreateTimeSlot(timeSlot)

	if err != nil {
		return ts, err
	}

	its := models.InterviewerTimeSlot{
		InterviewerId: iid2,
		TimeSlotId:    ts.Id,
	}

	err = db.Insert(&its)

	if err != nil {
		return ts, err
	}

	return ts, nil
}

func DeleteInterviewerTimeSlot(iid string, timeSlot *models.TimeSlot) error {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()
	iid2, err := strconv.Atoi(iid)

	if err != nil {
		return err
	}

	its := models.InterviewerTimeSlot{
		InterviewerId: iid2,
		TimeSlotId:    timeSlot.Id,
	}

	err = db.Model(&its).Select()

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	err = db.Delete(&its)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	err = DeleteTimeSlot(timeSlot)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return err
}
