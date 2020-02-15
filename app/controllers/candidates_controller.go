package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

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
	candidate := queries.FindCandidate(id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(candidate); err != nil {
		log.Fatal(err)
	}
}
