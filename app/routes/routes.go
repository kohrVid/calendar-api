package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kohrVid/calendar-api/app/controllers"
)

func Load() http.Handler {
	return routes()
}

func routes() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	candidatesAvailabilityResources(r)
	candidatesResources(r)
	interviewersAvailabilityResources(r)
	interviewersResources(r)
	timeSlotsResources(r)
	r.HandleFunc("/health", controllers.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", controllers.RootHandler).Methods("GET")
	return r
}

func candidatesResources(r *mux.Router) *mux.Router {
	r.HandleFunc(
		"/candidates/{id}",
		controllers.DeleteCandidatesHandler,
	).Methods("DELETE")

	r.HandleFunc(
		"/candidates/{id}",
		controllers.EditCandidatesHandler,
	).Methods("PATCH")

	r.HandleFunc(
		"/candidates/{id}",
		controllers.ShowCandidatesHandler,
	).Methods("GET")

	r.HandleFunc(
		"/candidates",
		controllers.NewCandidatesHandler,
	).Methods("POST")

	r.HandleFunc(
		"/candidates",
		controllers.CandidatesIndexHandler,
	).Methods("GET")

	return r
}

func candidatesAvailabilityResources(r *mux.Router) *mux.Router {
	r.HandleFunc(
		"/candidates/{cid}/availability/{id}",
		controllers.DeleteCandidateAvailabilityHandler,
	).Methods("DELETE")

	r.HandleFunc(
		"/candidates/{cid}/availability/{id}",
		controllers.EditCandidateAvailabilityHandler,
	).Methods("PATCH")

	r.HandleFunc(
		"/candidates/{cid}/availability/{id}",
		controllers.ShowCandidateAvailabilityHandler,
	).Methods("GET")

	r.HandleFunc(
		"/candidates/{cid}/availability",
		controllers.NewCandidateAvailabilityHandler,
	).Methods("POST")

	r.HandleFunc(
		"/candidates/{cid}/availability",
		controllers.CandidateAvailabilityIndexHandler,
	).Methods("GET")

	return r
}

func interviewersAvailabilityResources(r *mux.Router) *mux.Router {
	r.HandleFunc(
		"/interviewers/{iid}/availability/{id}",
		controllers.DeleteInterviewerAvailabilityHandler,
	).Methods("DELETE")

	r.HandleFunc(
		"/interviewers/{iid}/availability/{id}",
		controllers.EditInterviewerAvailabilityHandler,
	).Methods("PATCH")

	r.HandleFunc(
		"/interviewers/{iid}/availability/{id}",
		controllers.ShowInterviewerAvailabilityHandler,
	).Methods("GET")

	r.HandleFunc(
		"/interviewers/{iid}/availability",
		controllers.NewInterviewerAvailabilityHandler,
	).Methods("POST")

	r.HandleFunc(
		"/interviewers/{iid}/availability",
		controllers.InterviewerAvailabilityIndexHandler,
	).Methods("GET")

	return r
}

func interviewersResources(r *mux.Router) *mux.Router {
	r.HandleFunc(
		"/interviewers/{id}",
		controllers.DeleteInterviewersHandler,
	).Methods("DELETE")

	r.HandleFunc(
		"/interviewers/{id}",
		controllers.EditInterviewersHandler,
	).Methods("PATCH")

	r.HandleFunc(
		"/interviewers/{id}",
		controllers.ShowInterviewersHandler,
	).Methods("GET")

	r.HandleFunc(
		"/interviewers",
		controllers.NewInterviewersHandler,
	).Methods("POST")

	r.HandleFunc(
		"/interviewers",
		controllers.InterviewersIndexHandler,
	).Methods("GET")

	return r
}

func timeSlotsResources(r *mux.Router) *mux.Router {
	r.HandleFunc(
		"/time_slots/{id}",
		controllers.DeleteTimeSlotsHandler,
	).Methods("DELETE")

	r.HandleFunc(
		"/time_slots/{id}",
		controllers.EditTimeSlotsHandler,
	).Methods("PATCH")

	r.HandleFunc(
		"/time_slots/{id}",
		controllers.ShowTimeSlotsHandler,
	).Methods("GET")

	r.HandleFunc(
		"/time_slots",
		controllers.NewTimeSlotsHandler,
	).Methods("POST")

	r.HandleFunc(
		"/time_slots",
		controllers.TimeSlotsIndexHandler,
	).Methods("GET")
	r.HandleFunc(
		"/availability",
		controllers.TimeSlotsIndexHandler,
	).Methods("GET")

	return r
}
