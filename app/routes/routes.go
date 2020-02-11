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
	r.HandleFunc("/candidates", controllers.CandidatesHandler)
	r.HandleFunc("/health", controllers.HealthCheckHandler)
	r.HandleFunc("/", controllers.RootHandler)
	return r
}
