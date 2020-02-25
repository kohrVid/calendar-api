package commands

import (
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/kohrVid/calendar-api/db/sql/queries"
	"github.com/stretchr/testify/assert"
)

func init() {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestCreateTimeSlot(t *testing.T) {
	timeSlot := models.TimeSlot{
		Date:      "2020-02-25",
		StartTime: 13,
		Duration:  4,
	}

	res, err := CreateTimeSlot(&timeSlot)

	expected := models.TimeSlot{
		Id:        3,
		Date:      "2020-02-25",
		StartTime: timeSlot.StartTime,
		Duration:  timeSlot.Duration,
	}

	assert.Equal(t, expected, res, "New timeSlot expected")
	assert.Equal(t, nil, err, "No error expected")
}

func TestCreateTimeSlotWithMissingFields(t *testing.T) {
	timeSlot := models.TimeSlot{
		Date:      "2020-02-25",
		StartTime: 0,
		Duration:  0,
	}

	res, err := CreateTimeSlot(&timeSlot)

	assert.Equal(
		t,
		"ERROR #23502 null value in column \"start_time\" violates not-null constraint",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.Date, "No timeSlot details expected")
	assert.Equal(t, 0, res.StartTime, "No timeSlot details expected")
	assert.Equal(t, 0, res.Duration, "No timeSlot details expected")
}

func TestUpdateTimeSlot(t *testing.T) {
	timeSlot, err := queries.FindTimeSlot("1")

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	params := models.TimeSlot{
		Duration: 4,
	}

	res := UpdateTimeSlot(&timeSlot, params)

	updated_timeSlot, err := queries.FindTimeSlot("1")

	expected := models.TimeSlot{
		Id:        timeSlot.Id,
		Date:      "2020-02-25",
		StartTime: timeSlot.StartTime,
		Duration:  params.Duration,
	}

	assert.Equal(t, expected, res, "Updated timeSlot expected")
	assert.Equal(t, expected, updated_timeSlot, "Updated timeSlot expected")
}

func TestDeleteTimeSlot(t *testing.T) {
	id := "1"
	timeSlot, err := queries.FindTimeSlot(id)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	err = DeleteTimeSlot(&timeSlot)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	res, err := queries.FindTimeSlot(id)

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, 0, res.StartTime, "No timeSlot details expected")
	assert.Equal(t, 0, res.Duration, "No timeSlot details expected")
}
