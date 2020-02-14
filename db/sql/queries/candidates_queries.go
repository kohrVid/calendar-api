package queries

import (
	"fmt"
	"strconv"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
)

func ListCandidates() []models.Candidate {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	var candidates []models.Candidate
	err := db.Model(&candidates).Select()

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return candidates
}

func FindCandidate(id string) models.Candidate {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	idx, err := strconv.Atoi(id)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	candidate := &models.Candidate{Id: int64(idx)}
	err = db.Select(candidate)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return *candidate
}
