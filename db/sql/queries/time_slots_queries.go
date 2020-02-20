package queries

import (
	"fmt"
	"strconv"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
)

func ListTimeSlots() []models.TimeSlot {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)

	time_slots := make([]models.TimeSlot, 0)
	err := db.Model(&time_slots).Select()

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return time_slots
}

func FindTimeSlot(id string) (models.TimeSlot, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	idx, err := strconv.Atoi(id)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	time_slot := &models.TimeSlot{Id: idx}
	err = db.Select(time_slot)

	if err != nil {
		return *time_slot, err
	}

	return *time_slot, nil
}
