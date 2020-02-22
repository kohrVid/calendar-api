package queries

import (
	"log"
	"os"
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
	ret := m.Run()
	os.Exit(ret)
}

func TestListCandidates(t *testing.T) {
	conf := config.LoadConfig()

	users := config.ToMapList(
		conf["data"].(map[string]interface{})["candidates"],
	)

	candidate1 := models.Candidate{
		Id:        1,
		FirstName: users[0]["first_name"].(string),
		LastName:  users[0]["last_name"].(string),
		Email:     users[0]["email"].(string),
	}

	candidate2 := models.Candidate{
		Id:        2,
		FirstName: users[1]["first_name"].(string),
		LastName:  users[1]["last_name"].(string),
		Email:     users[1]["email"].(string),
	}

	expected := []models.Candidate{candidate1, candidate2}
	res := ListCandidates()

	assert.Equal(t, expected, res, "List of candidates expected")
}

func TestListCandidatesEmptyDB(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	res := ListCandidates()
	expected := []models.Candidate{}

	assert.Equal(t, expected, res, "Empty list expected")
	dbHelpers.Seed(conf)
}

func TestFindCandidate(t *testing.T) {
	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["candidates"],
	)[0]

	res, err := FindCandidate("1")

	if err != nil {
		log.Println(err)
	}

	expected := models.Candidate{
		Id:        1,
		FirstName: user["first_name"].(string),
		LastName:  user["last_name"].(string),
		Email:     user["email"].(string),
	}

	assert.Equal(t, expected, res, "first candidate expected")
	assert.Equal(t, nil, err, "No error expected")
}

func TestFindCandidateThatDoesNotExist(t *testing.T) {
	res, err := FindCandidate("10000")

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.FirstName, "No candidate details expected")
	assert.Equal(t, "", res.LastName, "No candidate details expected")
	assert.Equal(t, "", res.Email, "No candidate details expected")
}

func TestListCandidateTimeSlots(t *testing.T) {
	conf := config.LoadConfig()
	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[0]

	res := ListCandidateTimeSlots("1")

	timeSlot := models.TimeSlot{
		Id:        1,
		StartTime: ts["start_time"].(int),
		Duration:  ts["duration"].(int),
	}

	expected := []models.TimeSlot{timeSlot}

	assert.Equal(t, expected, res, "List of candidates expected")
}

func TestListCandidateTimeSlotsWhenEmpty(t *testing.T) {
	res := ListCandidateTimeSlots("2")
	expected := []models.TimeSlot{}

	assert.Equal(t, expected, res, "Empty list expected")
}

func TestFindCandidateTimeSlot(t *testing.T) {
	conf := config.LoadConfig()
	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[0]

	res, err := FindCandidateTimeSlot("1", "1")

	if err != nil {
		log.Println(err)
	}

	expected := models.TimeSlot{
		Id:        1,
		StartTime: ts["start_time"].(int),
		Duration:  ts["duration"].(int),
	}

	assert.Equal(t, expected, res, "List of candidates expected")
}

func TestFindCandidateTimeSlotDoesNotExist(t *testing.T) {
	res, err := FindCandidateTimeSlot("1", "2")

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
	assert.Equal(t, 0, res.Duration, "No time slot details expected")
}

func TestFindCandidateTimeSlotIsForAnotherCandidate(t *testing.T) {
	res, err := FindCandidateTimeSlot("2", "1")

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, 0, res.StartTime, "No time slot details expected")
	assert.Equal(t, 0, res.Duration, "No time slot details expected")
}
