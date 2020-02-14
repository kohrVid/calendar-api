package queries

import (
	"fmt"

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
