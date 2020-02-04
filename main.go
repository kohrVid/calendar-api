package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kohrVid/calendar-api/app/routes"
)

func main() {
	port := fmt.Sprintf(":%v", os.Getenv("PORT"))
	if port == ":" {
		port = ":8080"
	}

	http.ListenAndServe(port, routes.Load())
}
