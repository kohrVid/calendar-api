package models

import (
	"testing"

	"github.com/kohrVid/calendar-api/app/serializers"
	"github.com/kohrVid/calendar-api/config"
	"github.com/stretchr/testify/assert"
)

func TestToSerializer(t *testing.T) {
	conf := config.LoadConfig()
	user := config.ToMapList(conf["users"])[0]

	expected := serializers.Candidate{
		Id:        1,
		FirstName: user["first_name"].(string),
		LastName:  user["last_name"].(string),
		Email:     user["email"].(string),
	}

	candidate := Candidate{
		Id:        1,
		FirstName: user["first_name"].(string),
		LastName:  user["last_name"].(string),
		Email:     user["email"].(string),
	}

	assert.Equal(
		t,
		expected,
		candidate.ToSerializer(),
		"Expected to serialise the candidate",
	)
}
