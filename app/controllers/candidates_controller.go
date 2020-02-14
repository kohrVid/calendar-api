package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/kohrVid/calendar-api/app/serializers"
	"github.com/kohrVid/calendar-api/db/sql/queries"
)

func CandidatesHandler(w http.ResponseWriter, r *http.Request) {
	allCandidates := candidatesListPresenter()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(allCandidates); err != nil {
		log.Fatal(err)
	}
}

func candidatesListPresenter() []serializers.Candidate {
	cs := []serializers.Candidate{}
	allCandidates := queries.ListCandidates()

	for _, c := range allCandidates {
		cs = append(cs, c.ToSerializer())
	}

	return cs
}
