package controllers

import (
	"github.com/gorilla/mux"
)

func MockRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/candidates", CandidatesHandler).Methods("GET")
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", RootHandler).Methods("GET")
	return r
}
