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
	candidatesResources(r)
	r.HandleFunc("/health", controllers.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", controllers.RootHandler).Methods("GET")
	return r
}

func candidatesResources(r *mux.Router) *mux.Router {
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
		controllers.CandidatesIndexHandler,
	).Methods("GET")

	r.HandleFunc(
		"/candidates",
		controllers.NewCandidatesHandler,
	).Methods("POST")

	return r
}
