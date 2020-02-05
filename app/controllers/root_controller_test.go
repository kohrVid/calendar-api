package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func MockRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/", RootHandler).Methods("GET")
	return r
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "200 response expected")
	assert.Equal(t, "Calendar API", resp.Body.String(), "Page title expected")
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "200 response expected")
	assert.Equal(t, "OK", resp.Body.String(), "OK response is expected")
}