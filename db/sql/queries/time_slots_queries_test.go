package queries

import (
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/stretchr/testify/assert"
)

func TestListTimeSlots(t *testing.T) {
	conf := config.LoadConfig()

	timeSlots := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)

	res := ListTimeSlots()

	timeSlot1 := models.TimeSlot{
		Id:        1,
		Date:      timeSlots[0]["date"].(string),
		StartTime: timeSlots[0]["start_time"].(int),
		EndTime:   timeSlots[0]["end_time"].(int),
	}

	timeSlot2 := models.TimeSlot{
		Id:        2,
		Date:      timeSlots[1]["date"].(string),
		StartTime: timeSlots[1]["start_time"].(int),
		EndTime:   timeSlots[1]["end_time"].(int),
	}

	expected := []models.TimeSlot{timeSlot1, timeSlot2}

	assert.Equal(t, expected, res, "List of time slots expected")
}

func TestListTimeSlotsEmptyDB(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	res := ListTimeSlots()
	expected := []models.TimeSlot{}

	assert.Equal(t, expected, res, "Empty list expected")
	dbHelpers.Seed(conf)
}

func TestFindTimeSlot(t *testing.T) {
	conf := config.LoadConfig()

	timeSlot := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[0]

	res, err := FindTimeSlot("1")

	expected := models.TimeSlot{
		Id:        1,
		Date:      timeSlot["date"].(string),
		StartTime: timeSlot["start_time"].(int),
		EndTime:   timeSlot["end_time"].(int),
	}

	assert.Equal(t, expected, res, "first time slot expected")
	assert.Equal(t, nil, err, "No error expected")
}

func TestFindTimeSlotThatDoesNotExist(t *testing.T) {
	res, err := FindTimeSlot("10000")

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, 0, res.StartTime, "No time slot details expected")
	assert.Equal(t, 0, res.EndTime, "No time slot details expected")
}
