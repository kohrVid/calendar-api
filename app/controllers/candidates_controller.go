package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/kohrVid/calendar-api/db/sql/commands"
	"github.com/kohrVid/calendar-api/db/sql/queries"
)

func CandidatesIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	body := queries.ListCandidates()

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatal(err)
	}
}

func ShowCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
	candidate, err := queries.FindCandidate(id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(candidate); err != nil {
			log.Fatal(err)
		}
	}
}

func NewCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	candidate := new(models.Candidate)

	err = json.Unmarshal(c, candidate)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	_, err = commands.CreateCandidate(candidate)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		body := dbHelpers.PgErrorHandler(err, "candidates")
		fmt.Fprintf(w, body)
	} else {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(candidate)
	}
}

func EditCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]

	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	candidate, err := queries.FindCandidate(id)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	params := new(models.Candidate)

	err = json.Unmarshal(c, params)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	pa := *params

	commands.UpdateCandidate(&candidate, pa)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(candidate); err != nil {
		log.Fatal(err)
	}
}

func DeleteCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
	candidate, err := queries.FindCandidate(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		err = commands.DeleteCandidate(&candidate)

		if err != nil {
			fmt.Errorf("Error: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, fmt.Sprintf("Candidate #%v deleted", id))
	}
}

func CandidateAvailabilityIndexHandler(w http.ResponseWriter, r *http.Request) {
	cid := strings.Split(r.URL.Path, "/")[2]
	queryString := r.URL.Query()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	all := false
	var err error

	if (len(queryString) > 0) && (len(queryString["all"]) > 0) {
		all, err = strconv.ParseBool(queryString["all"][0])

		if err != nil {
			log.Fatal(err)
		}
	}

	if (len(queryString) > 0) && (len(queryString["interviewer"]) > 0) {
		interviewersIds := queryString["interviewer"]
		body := queries.ListCandidateAndInterviewerTimeSlot(cid, interviewersIds, all)

		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Fatal(err)
		}
	} else {
		body := queries.ListCandidateTimeSlots(cid, all)

		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Fatal(err)
		}
	}
}

func ShowCandidateAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	cid := path[2]
	id := path[4]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := queries.FindCandidateTimeSlot(cid, id)

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

func NewCandidateAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	cid := strings.Split(r.URL.Path, "/")[2]

	ts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	timeSlot := new(models.TimeSlot)

	err = json.Unmarshal(ts, timeSlot)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	_, err = commands.CreateCandidateTimeSlot(cid, timeSlot)
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

func EditCandidateAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	cid := path[2]
	id := path[4]

	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	timeSlot, err := queries.FindCandidateTimeSlot(cid, id)
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

func DeleteCandidateAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	cid := path[2]
	id := path[4]
	timeSlot, err := queries.FindCandidateTimeSlot(cid, id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(&models.Empty{}); err != nil {
			log.Fatal(err)
		}
	} else {
		err = commands.DeleteCandidateTimeSlot(cid, &timeSlot)

		if err != nil {
			fmt.Errorf("Error: %v", err)
			log.Printf("Error: %v\n", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, fmt.Sprintf("TimeSlot #%v deleted", id))
	}
}
