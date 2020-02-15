package queries

import (
	"os"
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
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
	users := config.ToMapList(conf["users"])

	db := db.DBConnect(conf)
	defer db.Close()

	res := ListCandidates()

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

	assert.Equal(t, expected, res, "List of candidates expected")
}

func TestFindCandidate(t *testing.T) {
	conf := config.LoadConfig()
	user := config.ToMapList(conf["users"])[0]

	db := db.DBConnect(conf)
	defer db.Close()

	res := FindCandidate("1")
	expected := models.Candidate{
		Id:        1,
		FirstName: user["first_name"].(string),
		LastName:  user["last_name"].(string),
		Email:     user["email"].(string),
	}

	assert.Equal(t, expected, res, "first candidate expected")
}
