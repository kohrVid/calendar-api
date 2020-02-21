package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/kohrVid/calendar-api/db/sql/commands"
	"github.com/kohrVid/calendar-api/db/sql/queries"
)

func TimeSlotsIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	body := queries.ListTimeSlots()

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatal(err)
	}
}

func ShowTimeSlotsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
	timeSlot, err := queries.FindTimeSlot(id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(timeSlot); err != nil {
			log.Fatal(err)
		}
	}
}

func NewTimeSlotsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	timeSlot := new(models.TimeSlot)

	err = json.Unmarshal(c, timeSlot)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	_, err = commands.CreateTimeSlot(timeSlot)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		body := dbHelpers.PgErrorHandler(err, "timeSlots")
		fmt.Fprintf(w, body)
	} else {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(timeSlot)
	}
}

func EditTimeSlotsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]

	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	timeSlot, err := queries.FindTimeSlot(id)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	params := new(models.TimeSlot)

	err = json.Unmarshal(c, params)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	pa := *params

	commands.UpdateTimeSlot(&timeSlot, pa)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(timeSlot); err != nil {
		log.Fatal(err)
	}
}

func DeleteTimeSlotsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
	timeSlot, err := queries.FindTimeSlot(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		err = commands.DeleteTimeSlot(&timeSlot)

		if err != nil {
			fmt.Errorf("Error: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, fmt.Sprintf("TimeSlot #%v deleted", id))
	}
}
