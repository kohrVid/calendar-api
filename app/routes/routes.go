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
	r := mux.NewRouter()
	r.HandleFunc("/candidates", controllers.CandidatesIndexHandler).Methods("GET")
	r.HandleFunc("/candidates/{id}", controllers.CandidatesShowHandler).Methods("GET")
	r.HandleFunc("/health", controllers.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", controllers.RootHandler).Methods("GET")
	return r
}
