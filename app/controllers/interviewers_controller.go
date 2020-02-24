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

func InterviewersIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	body := queries.ListInterviewers()

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatal(err)
	}
}

func ShowInterviewersHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
	interviewer, err := queries.FindInterviewer(id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(interviewer); err != nil {
			log.Fatal(err)
		}
	}
}

func NewInterviewersHandler(w http.ResponseWriter, r *http.Request) {
	i, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	interviewer := new(models.Interviewer)

	err = json.Unmarshal(i, interviewer)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	_, err = commands.CreateInterviewer(interviewer)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		body := dbHelpers.PgErrorHandler(err, "interviewers")
		fmt.Fprintf(w, body)
	} else {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(interviewer)
	}
}

func EditInterviewersHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]

	i, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	interviewer, err := queries.FindInterviewer(id)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	params := new(models.Interviewer)

	err = json.Unmarshal(i, params)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	pa := *params

	commands.UpdateInterviewer(&interviewer, pa)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(interviewer); err != nil {
		log.Fatal(err)
	}
}

func DeleteInterviewersHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
	interviewer, err := queries.FindInterviewer(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		err = commands.DeleteInterviewer(&interviewer)

		if err != nil {
			fmt.Errorf("Error: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, fmt.Sprintf("Interviewer #%v deleted", id))
	}
}

func InterviewerAvailabilityIndexHandler(w http.ResponseWriter, r *http.Request) {
	iid := strings.Split(r.URL.Path, "/")[2]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	body := queries.ListInterviewerTimeSlots(iid)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatal(err)
	}
}

func ShowInterviewerAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	iid := path[2]
	id := path[4]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := queries.FindInterviewerTimeSlot(iid, id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Fatal(err)
		}
	}
}

func NewInterviewerAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	iid := strings.Split(r.URL.Path, "/")[2]

	ts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	timeSlot := new(models.TimeSlot)

	err = json.Unmarshal(ts, timeSlot)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	_, err = commands.CreateInterviewerTimeSlot(iid, timeSlot)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		body := dbHelpers.PgErrorHandler(err, "time_slots")
		fmt.Fprintf(w, body)
	} else {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(timeSlot)
	}
}

func EditInterviewerAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	iid := path[2]
	id := path[4]

	i, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	timeSlot, err := queries.FindInterviewerTimeSlot(iid, id)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	params := new(models.TimeSlot)

	err = json.Unmarshal(i, params)
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

func DeleteInterviewerAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	iid := path[2]
	id := path[4]
	timeSlot, err := queries.FindInterviewerTimeSlot(iid, id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		err = commands.DeleteInterviewerTimeSlot(iid, &timeSlot)

		if err != nil {
			fmt.Errorf("Error: %v", err)
			log.Printf("Error: %v\n", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, fmt.Sprintf("TimeSlot #%v deleted", id))
	}
}
