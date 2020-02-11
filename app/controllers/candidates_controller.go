package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/kohrVid/calendar-api/app/serializers"
)

var candidates []serializers.Candidate = []serializers.Candidate{}

func CandidatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(candidates); err != nil {
		log.Fatal(err)
	}
}
