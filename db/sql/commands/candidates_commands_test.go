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
	dbHelpers.Clean(conf)
	os.Exit(ret)
}

func TestCreateCandidate(t *testing.T) {
	candidate := models.Candidate{
		FirstName: "Barnie",
		LastName:  "McAlister",
		Email:     "barnie.mcalister@example.com",
	}

	res := CreateCandidate(&candidate)

	expected := models.Candidate{
		Id:        3,
		FirstName: candidate.FirstName,
		LastName:  candidate.LastName,
		Email:     candidate.Email,
	}

	assert.Equal(t, expected, res, "New candidate expected")
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
