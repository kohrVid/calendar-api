package controllers

import (
	"github.com/gorilla/mux"
)

func MockRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/candidates", CandidatesIndexHandler).Methods("GET")
	r.HandleFunc("/candidates/{id}", CandidatesShowHandler).Methods("GET")
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", RootHandler).Methods("GET")
	return r
}
