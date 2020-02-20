package commands

import (
	"fmt"

	structs "github.com/fatih/structs"
	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
)

func CreateTimeSlot(timeSlot *models.TimeSlot) (models.TimeSlot, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	err := db.Insert(timeSlot)

	if err != nil {
		c := models.TimeSlot{}

		return c, err
	}

	c := *timeSlot

	return c, nil
}

func UpdateTimeSlot(timeSlot *models.TimeSlot, params models.TimeSlot) models.TimeSlot {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	t := structs.New(timeSlot)
	p := structs.New(params)

	sql := fmt.Sprintf(
		"UPDATE time_slots %v WHERE id = %v;",
		dbHelpers.SetSqlColumns(t, p),
		timeSlot.Id,
	)

	_, err := db.Model(t).Exec(sql)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	ts := *timeSlot

	return ts
}

func DeleteTimeSlot(timeSlot *models.TimeSlot) error {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	err := db.Delete(timeSlot)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return err
}
