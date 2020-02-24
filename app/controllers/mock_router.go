package controllers

import (
	"github.com/gorilla/mux"
)

func MockRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc(
		"/candidates/{cid}/availability/{id}",
		DeleteCandidateAvailabilityHandler,
	).Methods("DELETE")

	r.HandleFunc(
		"/candidates/{cid}/availability/{id}",
		EditCandidateAvailabilityHandler,
	).Methods("PATCH")

	r.HandleFunc(
		"/candidates/{cid}/availability/{id}",
		ShowCandidateAvailabilityHandler,
	).Methods("GET")
	r.HandleFunc(
		"/candidates/{cid}/availability",
		NewCandidateAvailabilityHandler,
	).Methods("POST")
	r.HandleFunc(
		"/candidates/{cid}/availability",
		CandidateAvailabilityIndexHandler,
	).Methods("GET")

	r.HandleFunc("/candidates/{id}", DeleteCandidatesHandler).Methods("DELETE")
	r.HandleFunc("/candidates/{id}", EditCandidatesHandler).Methods("PATCH")
	r.HandleFunc("/candidates/{id}", ShowCandidatesHandler).Methods("GET")
	r.HandleFunc("/candidates", NewCandidatesHandler).Methods("POST")
	r.HandleFunc("/candidates", CandidatesIndexHandler).Methods("GET")

	r.HandleFunc(
		"/interviewers/{iid}/availability/{id}",
		DeleteInterviewerAvailabilityHandler,
	).Methods("DELETE")

	r.HandleFunc(
		"/interviewers/{iid}/availability/{id}",
		EditInterviewerAvailabilityHandler,
	).Methods("PATCH")

	r.HandleFunc(
		"/interviewers/{iid}/availability/{id}",
		ShowInterviewerAvailabilityHandler,
	).Methods("GET")
	r.HandleFunc(
		"/interviewers/{iid}/availability",
		NewInterviewerAvailabilityHandler,
	).Methods("POST")
	r.HandleFunc(
		"/interviewers/{iid}/availability",
		InterviewerAvailabilityIndexHandler,
	).Methods("GET")

	r.HandleFunc("/interviewers/{id}", DeleteInterviewersHandler).Methods("DELETE")
	r.HandleFunc("/interviewers/{id}", EditInterviewersHandler).Methods("PATCH")
	r.HandleFunc("/interviewers/{id}", ShowInterviewersHandler).Methods("GET")
	r.HandleFunc("/interviewers", NewInterviewersHandler).Methods("POST")
	r.HandleFunc("/interviewers", InterviewersIndexHandler).Methods("GET")

	r.HandleFunc("/time_slots/{id}", DeleteTimeSlotsHandler).Methods("DELETE")
	r.HandleFunc("/time_slots/{id}", EditTimeSlotsHandler).Methods("PATCH")
	r.HandleFunc("/time_slots/{id}", ShowTimeSlotsHandler).Methods("GET")
	r.HandleFunc("/time_slots", NewTimeSlotsHandler).Methods("POST")
	r.HandleFunc("/time_slots", TimeSlotsIndexHandler).Methods("GET")

	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", RootHandler).Methods("GET")
	return r
}
