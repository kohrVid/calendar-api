package queries

import (
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
	"github.com/stretchr/testify/assert"
)

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

	assert.Equal(t, expected, res, "200 response expected")
}
