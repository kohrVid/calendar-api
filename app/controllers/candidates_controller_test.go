package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
	ret := m.Run()
	os.Exit(ret)
}

func TestCandidatesIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/candidates", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	conf := config.LoadConfig()

	users := config.ToMapList(
		conf["data"].(map[string]interface{})["candidates"],
	)

	user1 := users[0]
	user2 := users[1]

	expectedBody := fmt.Sprintf(
		`[{"id":1,"first_name":"%v","last_name":"%v","email":"%v"},{"id":2,"first_name":"%v","last_name":"%v","email":"%v"}]
`,
		user1["first_name"].(string),
		user1["last_name"].(string),
		user1["email"].(string),
		user2["first_name"].(string),
		user2["last_name"].(string),
		user2["email"].(string),
	)

	assert.Equal(t, 200, resp.Code, "200 response expected")

	assert.Equal(
		t,
		"application/json; charset=UTF-8",
		resp.Header().Get("Content-Type"),
		"JSON response expected",
	)

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"List of candidates expected",
	)
}

func TestCandidatesIndexHandlerEmpty(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)

	req, err := http.NewRequest("GET", "/candidates", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)
	expectedBody := fmt.Sprintf("[]\n")

	assert.Equal(t, 200, resp.Code, "200 response expected")

	assert.Equal(
		t,
		"application/json; charset=UTF-8",
		resp.Header().Get("Content-Type"),
		"JSON response expected",
	)

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"List of candidates expected",
	)

	dbHelpers.Seed(conf)
}

func TestShowCandidatesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/candidates/1", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["candidates"],
	)[0]

	expectedBody := fmt.Sprintf(
		`{"id":1,"first_name":"%v","last_name":"%v","email":"%v"}
`,
		user["first_name"].(string),
		user["last_name"].(string),
		user["email"].(string),
	)

	assert.Equal(t, 200, resp.Code, "200 response expected")

	assert.Equal(
		t,
		"application/json; charset=UTF-8",
		resp.Header().Get("Content-Type"),
		"JSON response expected",
	)

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"JSON of candidate expected",
	)
}

func TestShowCandidatesHandlerWhenCandidateDoesNotExist(t *testing.T) {
	req, err := http.NewRequest("GET", "/candidates/1000", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)
	expectedBody := "{}\n"

	assert.Equal(t, 404, resp.Code, "404 response expected")

	assert.Equal(
		t,
		"application/json; charset=UTF-8",
		resp.Header().Get("Content-Type"),
		"JSON response expected",
	)

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"Empty hash expected",
	)
}

func TestNewCandidatesHandler(t *testing.T) {
	user := models.Candidate{
		FirstName: "Barnie",
		LastName:  "McAlister",
		Email:     "barnie.mcalister@example.com",
	}

	data := []byte(
		fmt.Sprintf(
			`{"first_name": "%v", "last_name": "%v", "email": "%v"}`,
			user.FirstName,
			user.LastName,
			user.Email,
		),
	)

	req, err := http.NewRequest("POST", "/candidates", bytes.NewBuffer(data))

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := fmt.Sprintf(
		`{"id":3,"first_name":"%v","last_name":"%v","email":"%v"}
`,
		user.FirstName,
		user.LastName,
		user.Email,
	)

	assert.Equal(t, 201, resp.Code, "201 response expected")

	assert.Equal(
		t,
		"application/json; charset=UTF-8",
		resp.Header().Get("Content-Type"),
		"JSON response expected",
	)

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"New candidate expected",
	)
}

func TestNewCandidatesHandlerWhereAlreadyExists(t *testing.T) {
	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["candidates"],
	)[0]

	data := []byte(
		fmt.Sprintf(
			`{"first_name": "%v", "last_name": "%v", "email": "%v"}`,
			user["first_name"].(string),
			user["last_name"].(string),
			user["email"].(string),
		),
	)

	req, err := http.NewRequest("POST", "/candidates", bytes.NewBuffer(data))

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	assert.Equal(t, 304, resp.Code, "304 response expected")

	assert.Equal(
		t,
		"Candidate already exists",
		resp.Body.String(),
		"Duplicate message expected",
	)
}

func TestNewCandidatesHandlerMissingFields(t *testing.T) {
	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["candidates"],
	)[0]

	data := []byte(
		fmt.Sprintf(
			`{"first_name": "%v", "last_name": "%v", "email": "%v"}`,
			user["first_name"].(string),
			"",
			user["email"].(string),
		),
	)

	req, err := http.NewRequest("POST", "/candidates", bytes.NewBuffer(data))

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	assert.Equal(t, 304, resp.Code, "304 response expected")

	assert.Equal(
		t,
		"Missing field \"last_name\" in candidate",
		resp.Body.String(),
		"Missing field error expected",
	)
}

func TestEditCandidatesHandler(t *testing.T) {
	conf := config.LoadConfig()

	originalUser := config.ToMapList(
		conf["data"].(map[string]interface{})["candidates"],
	)[0]

	user := models.Candidate{
		Id:        1,
		FirstName: "Alexandra",
	}

	data := []byte(
		fmt.Sprintf(
			`{"first_name": "%v", "last_name": "%v", "email": "%v"}`,
			user.FirstName,
			originalUser["last_name"],
			originalUser["email"],
		),
	)

	req, err := http.NewRequest("PATCH", "/candidates/1", bytes.NewBuffer(data))

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := fmt.Sprintf(
		`{"id":1,"first_name":"%v","last_name":"%v","email":"%v"}
`,
		user.FirstName,
		originalUser["last_name"],
		originalUser["email"],
	)

	assert.Equal(t, 200, resp.Code, "200 response expected")

	assert.Equal(
		t,
		"application/json; charset=UTF-8",
		resp.Header().Get("Content-Type"),
		"JSON response expected",
	)

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"Updated candidate expected",
	)
}

func TestDeleteCandidatesHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/candidates/1", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := fmt.Sprintf("Candidate #%v deleted", 1)

	assert.Equal(t, 200, resp.Code, "200 response expected")

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"Deletion message expected",
	)
}

func TestDeleteCandidatesHandlerWhenCandidateDoesNotExist(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/candidates/1000", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)
	expectedBody := "{}\n"

	assert.Equal(t, 404, resp.Code, "404 response expected")

	assert.Equal(
		t,
		"application/json; charset=UTF-8",
		resp.Header().Get("Content-Type"),
		"JSON response expected",
	)

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"Empty hash expected",
	)
}

func TestNewCandidateAvailabilityHandler(t *testing.T) {
	conf := config.LoadConfig()

	timeSlot := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[0]

	data := []byte(
		fmt.Sprintf(
			`{"start_time":%v,"duration":%v}`,
			timeSlot["start_time"].(int),
			timeSlot["duration"].(int),
		),
	)

	// StrictSlash(true) doesn't work in the test environment for some reason
	req, err := http.NewRequest(
		"POST",
		"/candidates/1/availability/",
		bytes.NewBuffer(data),
	)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := fmt.Sprintf(
		`{"id":3,"start_time":%v,"duration":%v}
`,
		timeSlot["start_time"].(int),
		timeSlot["duration"].(int),
	)

	assert.Equal(t, 201, resp.Code, "201 response expected")

	assert.Equal(
		t,
		"application/json; charset=UTF-8",
		resp.Header().Get("Content-Type"),
		"JSON response expected",
	)

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"New timeSlot expected",
	)
}
