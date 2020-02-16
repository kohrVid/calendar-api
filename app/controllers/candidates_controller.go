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

func CandidatesIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(queries.ListCandidates()); err != nil {
		log.Fatal(err)
	}
}

func ShowCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	id := p[len(p)-1]
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
		fmt.Fprintf(w, dbHelpers.PgErrorHandler(err, "candidates"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(candidate)
}

func EditCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	id := p[len(p)-1]

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
	p := strings.Split(r.URL.Path, "/")
	id := p[len(p)-1]
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
