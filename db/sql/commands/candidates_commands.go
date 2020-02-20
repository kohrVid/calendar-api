package commands

import (
	"fmt"

	structs "github.com/fatih/structs"
	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
)

func CreateCandidate(candidate *models.Candidate) (models.Candidate, error) {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	err := db.Insert(candidate)

	if err != nil {
		c := models.Candidate{}

		return c, err
	}

	c := *candidate

	return c, nil
}

func UpdateCandidate(candidate *models.Candidate, params models.Candidate) models.Candidate {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	c := structs.New(candidate)
	p := structs.New(params)

	sql := fmt.Sprintf(
		"UPDATE candidates %v WHERE id = %v;",
		dbHelpers.SetSqlColumns(c, p),
		candidate.Id,
	)

	_, err := db.Model(c).Exec(sql)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	cc := *candidate

	return cc
}

func DeleteCandidate(candidate *models.Candidate) error {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	err := db.Delete(candidate)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	return err
}
