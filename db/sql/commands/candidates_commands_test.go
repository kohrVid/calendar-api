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

func TestCreateCandidate(t *testing.T) {
	candidate := models.Candidate{
		FirstName: "Barnie",
		LastName:  "McAlister",
		Email:     "barnie.mcalister@example.com",
	}

	res, err := CreateCandidate(&candidate)

	expected := models.Candidate{
		Id:        3,
		FirstName: candidate.FirstName,
		LastName:  candidate.LastName,
		Email:     candidate.Email,
	}

	assert.Equal(t, expected, res, "New candidate expected")
	assert.Equal(t, nil, err, "No error expected")
}

func TestCreateCandidateAlreadyExists(t *testing.T) {
	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["candidates"],
	)[0]

	candidate := models.Candidate{
		FirstName: user["first_name"].(string),
		LastName:  user["last_name"].(string),
		Email:     user["email"].(string),
	}

	res, err := CreateCandidate(&candidate)

	assert.Equal(
		t,
		"ERROR #23505 duplicate key value violates unique constraint \"candidates_unique_idx\"",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.FirstName, "No candidate details expected")
	assert.Equal(t, "", res.LastName, "No candidate details expected")
	assert.Equal(t, "", res.Email, "No candidate details expected")
}

func TestCreateCandidateTimeSlotWithMissingFields(t *testing.T) {
	candidate := models.Candidate{
		FirstName: "",
		LastName:  "",
		Email:     "",
	}

	res, err := CreateCandidate(&candidate)

	assert.Equal(
		t,
		"ERROR #23502 null value in column \"first_name\" violates not-null constraint",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.FirstName, "No candidate details expected")
	assert.Equal(t, "", res.LastName, "No candidate details expected")
	assert.Equal(t, "", res.Email, "No candidate details expected")
}

func TestUpdateCandidate(t *testing.T) {
	candidate, err := queries.FindCandidate("1")

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	params := models.Candidate{
		FirstName: "Alexandra",
	}

	res := UpdateCandidate(&candidate, params)

	updatedCandidate, err := queries.FindCandidate("1")

	expected := models.Candidate{
		Id:        candidate.Id,
		FirstName: params.FirstName,
		LastName:  candidate.LastName,
		Email:     candidate.Email,
	}

	assert.Equal(t, expected, res, "Updated candidate expected")
	assert.Equal(t, expected, updatedCandidate, "Updated candidate expected")
}

func TestDeleteCandidate(t *testing.T) {
	conf := config.LoadConfig()
	id := "1"
	candidate, err := queries.FindCandidate(id)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	err = DeleteCandidate(&candidate)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	res, err := queries.FindCandidate(id)

	assert.Equal(
		t,
		"pg: no rows in result set",
		err.Error(),
		"Error expected",
	)

	assert.Equal(t, "", res.FirstName, "No candidate details expected")
	assert.Equal(t, "", res.LastName, "No candidate details expected")
	assert.Equal(t, "", res.Email, "No candidate details expected")

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestCreateCandidateTimeSlot(t *testing.T) {
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

	res, err := CreateCandidateTimeSlot("2", &timeSlot)

	expected := models.TimeSlot{
		Id:        3,
		Date:      timeSlot.Date,
		StartTime: timeSlot.StartTime,
		Duration:  timeSlot.Duration,
	}

	expectedCandidateTimeSlot := models.CandidateTimeSlot{
		Id:          2,
		CandidateId: 2,
		TimeSlotId:  3,
	}

	findCandidateTimeSlot := models.CandidateTimeSlot{Id: 2}

	err = db.Select(
		&findCandidateTimeSlot,
	)

	if err != nil {
		log.Println(err)
	}

	assert.Equal(t, expected, res, "New timeSlot expected")
	assert.Equal(t, nil, err, "No error expected")

	assert.Equal(
		t,
		expectedCandidateTimeSlot,
		findCandidateTimeSlot,
		"New candidate_time_slot association expected",
	)
}

func TestCreateCandidateTimeSlotIfAlreadyExists(t *testing.T) {
	conf := config.LoadConfig()

	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[0]

	timeSlot := models.TimeSlot{
		Date:      ts["date"].(string),
		StartTime: ts["start_time"].(int),
		Duration:  ts["duration"].(int),
	}

	res, err := CreateCandidateTimeSlot("1", &timeSlot)

	assert.Equal(
		t,
		"ERROR #23505 time slot already exists for candidate",
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

func TestCreateCandidateWithMissingFields(t *testing.T) {
	conf := config.LoadConfig()

	ts := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	timeSlot := models.TimeSlot{
		Duration: ts["duration"].(int),
	}

	res, err := CreateCandidateTimeSlot("2", &timeSlot)

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

func TestDeleteCandidateTimeSlot(t *testing.T) {
	conf := config.LoadConfig()
	cid := "1"
	id := "1"
	timeSlot, err := queries.FindCandidateTimeSlot(cid, id)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	err = DeleteCandidateTimeSlot(cid, &timeSlot)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	res, err := queries.FindCandidateTimeSlot(cid, id)

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
