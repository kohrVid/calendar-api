package commands

import (
	"log"
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/kohrVid/calendar-api/db/sql/queries"
	"github.com/stretchr/testify/assert"
)

func init() {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestCreateInterviewer(t *testing.T) {
	interviewer := models.Interviewer{
		FirstName: "Gwyneira",
		LastName:  "Vega",
		Email:     "gwyneira.vega@example.com",
	}

	res, err := CreateInterviewer(&interviewer)

	expected := models.Interviewer{
		Id:        3,
		FirstName: interviewer.FirstName,
		LastName:  interviewer.LastName,
		Email:     interviewer.Email,
	}

	assert.Equal(t, expected, res, "New interviewer expected")
	assert.Equal(t, nil, err, "No error expected")
}

func TestCreateInterviewerAlreadyExists(t *testing.T) {
	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["interviewers"],
	)[0]

	interviewer := models.Interviewer{
		FirstName: user["first_name"].(string),
		LastName:  user["last_name"].(string),
		Email:     user["email"].(string),
	}

	res, err := CreateInterviewer(&interviewer)

	assert.Equal(
		t,
		"ERROR #23505 duplicate key value violates unique constraint \"interviewers_unique_idx\"",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.FirstName, "No interviewer details expected")
	assert.Equal(t, "", res.LastName, "No interviewer details expected")
	assert.Equal(t, "", res.Email, "No interviewer details expected")
}

func TestCreateInterviewerTimeSlotWithMissingFields(t *testing.T) {
	interviewer := models.Interviewer{
		FirstName: "",
		LastName:  "",
		Email:     "",
	}

	res, err := CreateInterviewer(&interviewer)

	assert.Equal(
		t,
		"ERROR #23502 null value in column \"first_name\" violates not-null constraint",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.FirstName, "No interviewer details expected")
	assert.Equal(t, "", res.LastName, "No interviewer details expected")
	assert.Equal(t, "", res.Email, "No interviewer details expected")
}

func TestUpdateInterviewer(t *testing.T) {
	interviewer, err := queries.FindInterviewer("1")

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	params := models.Interviewer{
		FirstName: "Eira",
	}

	res := UpdateInterviewer(&interviewer, params)

	updatedInterviewer, err := queries.FindInterviewer("1")

	expected := models.Interviewer{
		Id:        interviewer.Id,
		FirstName: params.FirstName,
		LastName:  interviewer.LastName,
		Email:     interviewer.Email,
	}

	assert.Equal(t, expected, res, "Updated interviewer expected")
	assert.Equal(t, expected, updatedInterviewer, "Updated interviewer expected")
}

func TestDeleteInterviewer(t *testing.T) {
	conf := config.LoadConfig()
	id := "1"
	interviewer, err := queries.FindInterviewer(id)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	err = DeleteInterviewer(&interviewer)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	res, err := queries.FindInterviewer(id)

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.FirstName, "No interviewer details expected")
	assert.Equal(t, "", res.LastName, "No interviewer details expected")
	assert.Equal(t, "", res.Email, "No interviewer details expected")

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestCreateInterviewerTimeSlot(t *testing.T) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)

	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	timeSlot := models.TimeSlot{
		Date:      ts["date"].(string),
		StartTime: ts["start_time"].(int),
		Duration:  ts["duration"].(int),
	}

	res, err := CreateInterviewerTimeSlot("2", &timeSlot)

	expected := models.TimeSlot{
		Id:        3,
		Date:      timeSlot.Date,
		StartTime: timeSlot.StartTime,
		Duration:  timeSlot.Duration,
	}

	expectedInterviewerTimeSlot := models.InterviewerTimeSlot{
		Id:            2,
		InterviewerId: 2,
		TimeSlotId:    3,
	}

	findInterviewerTimeSlot := models.InterviewerTimeSlot{Id: 2}

	err = db.Select(
		&findInterviewerTimeSlot,
	)

	if err != nil {
		log.Println(err)
	}

	assert.Equal(t, expected, res, "New timeSlot expected")
	assert.Equal(t, nil, err, "No error expected")

	assert.Equal(
		t,
		expectedInterviewerTimeSlot,
		findInterviewerTimeSlot,
		"New interviewer_time_slot association expected",
	)
}

func TestCreateInterviewerTimeSlotIfAlreadyExists(t *testing.T) {
	conf := config.LoadConfig()

	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	timeSlot := models.TimeSlot{
		Date:      ts["date"].(string),
		StartTime: ts["start_time"].(int),
		Duration:  ts["duration"].(int),
	}

	res, err := CreateInterviewerTimeSlot("1", &timeSlot)

	assert.Equal(
		t,
		"ERROR #23505 time slot already exists for interviewer",
		err.Error(),
		"Error expected",
	)

	assert.Equal(
		t,
		timeSlot.Date,
		res.Date,
		"No time slot details expected",
	)

	assert.Equal(
		t,
		timeSlot.StartTime,
		res.StartTime,
		"No time slot details expected",
	)

	assert.Equal(
		t,
		timeSlot.Duration,
		res.Duration,
		"No time slot details expected",
	)
}

func TestCreateInterviewerWithMissingFields(t *testing.T) {
	conf := config.LoadConfig()

	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	timeSlot := models.TimeSlot{
		Duration: ts["duration"].(int),
	}

	res, err := CreateInterviewerTimeSlot("2", &timeSlot)

	assert.Equal(
		t,
		"ERROR #23502 null value in column \"start_time\" violates not-null constraint",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, 0, res.StartTime, "No time slot details expected")
	assert.Equal(t, 0, res.Duration, "No time slot details expected")

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestDeleteInterviewerTimeSlot(t *testing.T) {
	conf := config.LoadConfig()
	iid := "1"
	id := "2"
	timeSlot, err := queries.FindInterviewerTimeSlot(iid, id)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	err = DeleteInterviewerTimeSlot(iid, &timeSlot)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	res, err := queries.FindInterviewerTimeSlot(iid, id)

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, 0, res.StartTime, "No time slot details expected")
	assert.Equal(t, 0, res.Duration, "No time slot details expected")

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}
