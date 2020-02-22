package controllers

import (
	"github.com/gorilla/mux"
)

func MockRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/candidates/{cid}/availability/{id}", ShowCandidateAvailabilityHandler).Methods("GET")
	r.HandleFunc("/candidates/{cid}/availability", CandidateAvailabilityIndexHandler).Methods("GET")

	r.HandleFunc("/candidates/{id}", DeleteCandidatesHandler).Methods("DELETE")
	r.HandleFunc("/candidates/{id}", EditCandidatesHandler).Methods("PATCH")
	r.HandleFunc("/candidates/{id}", ShowCandidatesHandler).Methods("GET")
	r.HandleFunc("/candidates", NewCandidatesHandler).Methods("POST")
	r.HandleFunc("/candidates", CandidatesIndexHandler).Methods("GET")

	r.HandleFunc("/time_slots/{id}", DeleteTimeSlotsHandler).Methods("DELETE")
	r.HandleFunc("/time_slots/{id}", EditTimeSlotsHandler).Methods("PATCH")
	r.HandleFunc("/time_slots/{id}", ShowTimeSlotsHandler).Methods("GET")
	r.HandleFunc("/time_slots", NewTimeSlotsHandler).Methods("POST")
	r.HandleFunc("/time_slots", TimeSlotsIndexHandler).Methods("GET")

	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", RootHandler).Methods("GET")
	return r
}
