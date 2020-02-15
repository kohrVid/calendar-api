package controllers

import (
	"github.com/gorilla/mux"
)

func MockRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/candidates/{id}", EditCandidatesHandler).Methods("PATCH")
	r.HandleFunc("/candidates/{id}", ShowCandidatesHandler).Methods("GET")
	r.HandleFunc("/candidates", NewCandidatesHandler).Methods("POST")
	r.HandleFunc("/candidates", CandidatesIndexHandler).Methods("GET")
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", RootHandler).Methods("GET")
	return r
}
