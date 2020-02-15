package commands

import (
	"context"
	"fmt"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
)

func CreateCandidate(candidate *models.Candidate) models.Candidate {
	conf := config.LoadConfig()
	ctx := context.Background()
	db := db.DBConnect(conf).WithContext(ctx)

	err := db.Insert(candidate)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	c := *candidate

	return c
}
