package queries

import (
	"log"
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/stretchr/testify/assert"
)

func init() {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestListInterviewers(t *testing.T) {
	conf := config.LoadConfig()

	users := config.ToMapList(
		conf["data"].(map[string]interface{})["interviewers"],
	)

	interviewer1 := models.Interviewer{
		Id:        1,
		FirstName: users[0]["first_name"].(string),
		LastName:  users[0]["last_name"].(string),
		Email:     users[0]["email"].(string),
	}

	interviewer2 := models.Interviewer{
		Id:        2,
		FirstName: users[1]["first_name"].(string),
		LastName:  users[1]["last_name"].(string),
		Email:     users[1]["email"].(string),
	}

	expected := []models.Interviewer{interviewer1, interviewer2}
	res := ListInterviewers()

	assert.Equal(t, expected, res, "List of interviewers expected")
}

func TestListInterviewersEmptyDB(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	res := ListInterviewers()
	expected := []models.Interviewer{}

	assert.Equal(t, expected, res, "Empty list expected")
	dbHelpers.Seed(conf)
}

func TestFindInterviewer(t *testing.T) {
	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["interviewers"],
	)[0]

	res, err := FindInterviewer("1")

	if err != nil {
		log.Println(err)
	}

	expected := models.Interviewer{
		Id:        1,
		FirstName: user["first_name"].(string),
		LastName:  user["last_name"].(string),
		Email:     user["email"].(string),
	}

	assert.Equal(t, expected, res, "first interviewer expected")
	assert.Equal(t, nil, err, "No error expected")
}

func TestFindInterviewerThatDoesNotExist(t *testing.T) {
	res, err := FindInterviewer("10000")

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.FirstName, "No interviewer details expected")
	assert.Equal(t, "", res.LastName, "No interviewer details expected")
	assert.Equal(t, "", res.Email, "No interviewer details expected")
}

func TestListInterviewerTimeSlots(t *testing.T) {
	conf := config.LoadConfig()
	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	res := ListInterviewerTimeSlots("1")

	timeSlot := models.TimeSlot{
		Id:        2,
		Date:      ts["date"].(string),
		StartTime: ts["start_time"].(int),
		EndTime:   ts["end_time"].(int),
	}

	expected := []models.TimeSlot{timeSlot}

	assert.Equal(t, expected, res, "List of interviewers expected")
}

func TestListInterviewerTimeSlotsWhenEmpty(t *testing.T) {
	res := ListInterviewerTimeSlots("2")
	expected := []models.TimeSlot{}

	assert.Equal(t, expected, res, "Empty list expected")
}

func TestFindInterviewerTimeSlot(t *testing.T) {
	conf := config.LoadConfig()
	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	res, err := FindInterviewerTimeSlot("1", "2")

	if err != nil {
		log.Println(err)
	}

	expected := models.TimeSlot{
		Id:        2,
		Date:      ts["date"].(string),
		StartTime: ts["start_time"].(int),
		EndTime:   ts["end_time"].(int),
	}

	assert.Equal(t, expected, res, "List of interviewers expected")
}

func TestFindInterviewerTimeSlotDoesNotExist(t *testing.T) {
	res, err := FindInterviewerTimeSlot("1", "3")

	if err != nil {
		log.Println(err)
	}

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, 0, res.StartTime, "No time slot details expected")
	assert.Equal(t, 0, res.EndTime, "No time slot details expected")
}

func TestFindInterviewerTimeSlotIsForAnotherInterviewer(t *testing.T) {
	res, err := FindInterviewerTimeSlot("2", "2")

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, 0, res.StartTime, "No time slot details expected")
	assert.Equal(t, 0, res.EndTime, "No time slot details expected")
}
