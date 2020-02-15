package commands

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
)

func CreateCandidate(candidate *models.Candidate) models.Candidate {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)

	err := db.Insert(candidate)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	c := *candidate

	return c
}

func UpdateCandidate(candidate *models.Candidate, params models.Candidate) models.Candidate {
	conf := config.LoadConfig()
	db := db.DBConnect(conf)
	defer db.Close()

	c := structs.New(candidate)
	p := structs.New(params)

	for _, k := range c.Names() {
		if !p.Field(k).IsZero() {
			c.Field(k).Set(p.Field(k).Value())
		}
	}

	_, err := db.Model(c).Update()
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
