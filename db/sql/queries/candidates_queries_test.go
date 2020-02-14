package queries

import (
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	conf := config.LoadConfig()
	user1 := conf["user1"].(map[string]interface{})
	user2 := conf["user2"].(map[string]interface{})

	db := db.DBConnect(conf)
	defer db.Close()

	res := ListCandidates()
	candidate1 := models.Candidate{
		Id:        1,
		FirstName: user1["first_name"].(string),
		LastName:  user1["last_name"].(string),
		Email:     user1["email"].(string),
	}

	candidate2 := models.Candidate{
		Id:        2,
		FirstName: user2["first_name"].(string),
		LastName:  user2["last_name"].(string),
		Email:     user2["email"].(string),
	}

	expected := []models.Candidate{candidate1, candidate2}

	assert.Equal(t, expected, res, "200 response expected")
}
