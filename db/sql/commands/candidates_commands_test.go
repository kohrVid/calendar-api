package commands

import (
	"os"
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/kohrVid/calendar-api/db/sql/queries"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
	ret := m.Run()
	os.Exit(ret)
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
	user := config.ToMapList(conf["users"])[0]

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

func TestCreateCandidateWithMissingFields(t *testing.T) {
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

	expected := models.Candidate{
		Id:        candidate.Id,
		FirstName: params.FirstName,
		LastName:  candidate.LastName,
		Email:     candidate.Email,
	}

	assert.Equal(t, expected, res, "Updated candidate expected")
}

func TestDeleteCandidate(t *testing.T) {
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
}
