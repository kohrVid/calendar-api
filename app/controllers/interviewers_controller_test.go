package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/stretchr/testify/assert"
)

func init() {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestInterviewersIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/interviewers", nil)

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
		conf["data"].(map[string]interface{})["interviewers"],
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
		"List of interviewers expected",
	)
}

func TestInterviewersIndexHandlerEmpty(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)

	req, err := http.NewRequest("GET", "/interviewers", nil)

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
		"List of interviewers expected",
	)

	dbHelpers.Seed(conf)
}

func TestShowInterviewersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/interviewers/1", nil)

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
		conf["data"].(map[string]interface{})["interviewers"],
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
		"JSON of interviewer expected",
	)
}

func TestShowInterviewersHandlerWhenInterviewerDoesNotExist(t *testing.T) {
	req, err := http.NewRequest("GET", "/interviewers/1000", nil)

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

func TestNewInterviewersHandler(t *testing.T) {
	user := models.Interviewer{
		FirstName: "Gwyneira",
		LastName:  "Vega",
		Email:     "gwyneira.vega@example.com",
	}

	data := []byte(
		fmt.Sprintf(
			`{"first_name": "%v", "last_name": "%v", "email": "%v"}`,
			user.FirstName,
			user.LastName,
			user.Email,
		),
	)

	req, err := http.NewRequest("POST", "/interviewers", bytes.NewBuffer(data))

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
		"New interviewer expected",
	)
}

func TestNewInterviewersHandlerWhereAlreadyExists(t *testing.T) {
	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["interviewers"],
	)[0]

	data := []byte(
		fmt.Sprintf(
			`{"first_name": "%v", "last_name": "%v", "email": "%v"}`,
			user["first_name"].(string),
			user["last_name"].(string),
			user["email"].(string),
		),
	)

	req, err := http.NewRequest("POST", "/interviewers", bytes.NewBuffer(data))

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
		"Interviewer already exists",
		resp.Body.String(),
		"Duplicate message expected",
	)
}

func TestNewInterviewersHandlerMissingFields(t *testing.T) {
	conf := config.LoadConfig()

	user := config.ToMapList(
		conf["data"].(map[string]interface{})["interviewers"],
	)[0]

	data := []byte(
		fmt.Sprintf(
			`{"first_name": "%v", "last_name": "%v", "email": "%v"}`,
			user["first_name"].(string),
			"",
			user["email"].(string),
		),
	)

	req, err := http.NewRequest("POST", "/interviewers", bytes.NewBuffer(data))

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
		"Missing field \"last_name\" in interviewer",
		resp.Body.String(),
		"Missing field error expected",
	)
}

func TestEditInterviewersHandler(t *testing.T) {
	conf := config.LoadConfig()

	originalUser := config.ToMapList(
		conf["data"].(map[string]interface{})["interviewers"],
	)[0]

	user := models.Interviewer{
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

	req, err := http.NewRequest("PATCH", "/interviewers/1", bytes.NewBuffer(data))

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
		"Updated interviewer expected",
	)
}

func TestDeleteInterviewersHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/interviewers/1", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := fmt.Sprintf("Interviewer #%v deleted", 1)

	assert.Equal(t, 200, resp.Code, "200 response expected")

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"Deletion message expected",
	)
}

func TestDeleteInterviewersHandlerWhenInterviewerDoesNotExist(t *testing.T) {
	conf := config.LoadConfig()
	req, err := http.NewRequest("DELETE", "/interviewers/1000", nil)

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

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestInterviewerAvailabilityIndexHandler(t *testing.T) {
	conf := config.LoadConfig()

	req, err := http.NewRequest("GET", "/interviewers/1/availability", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	timeSlot := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	expectedBody := fmt.Sprintf(
		`[{"id":2,"start_time":%v,"duration":%v}]
`,
		timeSlot["start_time"].(int),
		timeSlot["duration"].(int),
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
		"List of time slots expected",
	)
}

func TestShowInterviewerAvailabilityHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/interviewers/1/availability/2", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)
	conf := config.LoadConfig()

	timeSlot := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	expectedBody := fmt.Sprintf(
		`{"id":2,"start_time":%v,"duration":%v}
`,
		timeSlot["start_time"].(int),
		timeSlot["duration"].(int),
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
		"JSON of interviewer expected",
	)
}

func TestShowInterviewerAvailabilityHandlerWhenTimeslotDoesNotExist(t *testing.T) {
	req, err := http.NewRequest("GET", "/interviewers/1/availability/1", nil)

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

func TestShowInterviewerAvailabilityHandlerWhenTimeslotIsForDifferentInterviewer(t *testing.T) {
	req, err := http.NewRequest("GET", "/interviewers/2/availability/2", nil)

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

func TestNewInterviewerAvailabilityHandler(t *testing.T) {
	timeSlot := models.TimeSlot{
		StartTime: 13,
		Duration:  3,
	}

	data := []byte(
		fmt.Sprintf(
			`{"start_time":%v,"duration":%v}`,
			timeSlot.StartTime,
			timeSlot.Duration,
		),
	)

	req, err := http.NewRequest("POST", "/interviewers/2/availability", bytes.NewBuffer(data))

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
		timeSlot.StartTime,
		timeSlot.Duration,
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
		"New time slot expected",
	)
}

func TestNewInterviewerAvailabilityHandlerMissingFields(t *testing.T) {
	conf := config.LoadConfig()

	timeSlot := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[0]

	data := []byte(
		fmt.Sprintf(
			`{"start_time":%v}`,
			timeSlot["start_time"].(int),
		),
	)

	req, err := http.NewRequest("POST", "/interviewers/2/availability", bytes.NewBuffer(data))

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
		"Missing field \"duration\" in time_slot",
		resp.Body.String(),
		"Missing field error expected",
	)
}

func TestEditInterviewersAvailabilityHandler(t *testing.T) {
	conf := config.LoadConfig()

	originalTimeSlot := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[1]

	timeSlot := models.TimeSlot{
		Id:       1,
		Duration: 4,
	}

	data := []byte(
		fmt.Sprintf(
			`{"start_time":%v,"duration":%v}`,
			originalTimeSlot["start_time"].(int),
			timeSlot.Duration,
		),
	)

	req, err := http.NewRequest(
		"PATCH",
		"/interviewers/1/availability/2",
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
		`{"id":2,"start_time":%v,"duration":%v}
`,
		originalTimeSlot["start_time"].(int),
		timeSlot.Duration,
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
		"Updated time slot expected",
	)
}

func TestDeleteInterviewersAvailabilityHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/interviewers/1/availability/2", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := fmt.Sprintf("TimeSlot #%v deleted", 2)

	assert.Equal(t, 200, resp.Code, "200 response expected")

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"Deletion message expected",
	)
}

func TestDeleteInterviewersAvailabilityHandlerWhenInterviewerDoesNotExist(t *testing.T) {
	conf := config.LoadConfig()
	req, err := http.NewRequest("DELETE", "/interviewers/1000/availability/2", nil)

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

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}
func TestDeleteInterviewersAvailabilityHandlerWhenTimeSlotDoesNotExist(t *testing.T) {
	conf := config.LoadConfig()
	req, err := http.NewRequest("DELETE", "/interviewers/1/availability/1000", nil)

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

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}
