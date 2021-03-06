package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kohrVid/calendar-api/app/models"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/kohrVid/calendar-api/db/sql/commands"
	"github.com/stretchr/testify/assert"
)

var (
	pastTime    = "2020-02-24T07:00:48Z"
	currentTime = "2020-03-05T07:00:48Z"
)

func init() {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
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
	os.Setenv("TIME_NOW", fmt.Sprintf("'%v'::timestamp", pastTime))

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
	conf := config.LoadConfig()
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

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestCandidateAvailabilityIndexHandlerDefault(t *testing.T) {
	os.Setenv("TIME_NOW", fmt.Sprintf("'%v'::timestamp", currentTime))

	req, err := http.NewRequest("GET", "/candidates/1/availability", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := `[]
`

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

func TestCandidateAvailabilityIndexHandlerAllIsTrue(t *testing.T) {
	conf := config.LoadConfig()
	os.Setenv("TIME_NOW", fmt.Sprintf("'%v'::timestamp", currentTime))

	req, err := http.NewRequest("GET", "/candidates/1/availability?all=true", nil)

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
	)[0]

	expectedBody := fmt.Sprintf(
		`[{"id":1,"date":"%v","start_time":%v,"end_time":%v}]
`,
		timeSlot["date"].(string),
		timeSlot["start_time"].(int),
		timeSlot["end_time"].(int),
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

func TestCandidateAvailabilityIndexHandlerAllIsFalse(t *testing.T) {
	os.Setenv("TIME_NOW", fmt.Sprintf("'%v'::timestamp", currentTime))
	req, err := http.NewRequest("GET", "/candidates/1/availability?all=false", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := `[]
`

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

func TestCandidateAvailabilityIndexHandlerWithInterviewers(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
	os.Setenv("TIME_NOW", fmt.Sprintf("'%v'::timestamp", pastTime))

	expectedAvailability := models.TimeSlot{
		Id:        3,
		Date:      "2020-02-25",
		StartTime: 11,
		EndTime:   16,
	}

	_, err := commands.CreateInterviewerTimeSlot("2", &expectedAvailability)

	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(
		"GET",
		"/candidates/1/availability?interviewer=1&interviewer=2",
		nil,
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
		`[{"interviewer_id":1,"time_slots":[]},{"interviewer_id":2,"time_slots":[{"id":%v,"date":"%v","start_time":%v,"end_time":%v}]}]
`,
		expectedAvailability.Id,
		expectedAvailability.Date,
		expectedAvailability.StartTime,
		12,
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

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestCandidateAvailabilityIndexHandlerWithInterviewersAndAllIsTrue(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
	os.Setenv("TIME_NOW", fmt.Sprintf("'%v'::timestamp", currentTime))

	expectedAvailability := models.TimeSlot{
		Id:        3,
		Date:      "2020-02-25",
		StartTime: 11,
		EndTime:   16,
	}

	_, err := commands.CreateInterviewerTimeSlot("2", &expectedAvailability)

	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(
		"GET",
		"/candidates/1/availability?interviewer=1&interviewer=2&all=true",
		nil,
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
		`[{"interviewer_id":1,"time_slots":[]},{"interviewer_id":2,"time_slots":[{"id":%v,"date":"%v","start_time":%v,"end_time":%v}]}]
`,
		expectedAvailability.Id,
		expectedAvailability.Date,
		expectedAvailability.StartTime,
		12,
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

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestCandidateAvailabilityIndexHandlerWithInterviewersAndAllIsFalse(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
	os.Setenv("TIME_NOW", fmt.Sprintf("'%v'::timestamp", currentTime))

	expectedAvailability := models.TimeSlot{
		Id:        3,
		Date:      "2020-02-25",
		StartTime: 11,
		EndTime:   16,
	}

	_, err := commands.CreateInterviewerTimeSlot("2", &expectedAvailability)

	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(
		"GET",
		"/candidates/1/availability?interviewer=1&interviewer=2&all=false",
		nil,
	)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := `[{"interviewer_id":1,"time_slots":[]},{"interviewer_id":2,"time_slots":[]}]
`

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

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestCandidateAvailabilityIndexHandlerWithOnlyOneInterviewer(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
	os.Setenv("TIME_NOW", fmt.Sprintf("'%v'::timestamp", currentTime))

	expectedAvailability := models.TimeSlot{
		Id:        3,
		Date:      "2020-02-25",
		StartTime: 11,
		EndTime:   16,
	}

	_, err := commands.CreateInterviewerTimeSlot("2", &expectedAvailability)

	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(
		"GET",
		"/candidates/1/availability?interviewer=1",
		nil,
	)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := `[{"interviewer_id":1,"time_slots":[]}]
`

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

	dbHelpers.Clean(conf)
	dbHelpers.Seed(conf)
}

func TestShowCandidateAvailabilityHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/candidates/1/availability/1", nil)

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
	)[0]

	expectedBody := fmt.Sprintf(
		`{"id":1,"date":"%v","start_time":%v,"end_time":%v}
`,
		timeSlot["date"].(string),
		timeSlot["start_time"].(int),
		timeSlot["end_time"].(int),
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

func TestShowCandidateAvailabilityHandlerWhenTimeslotDoesNotExist(t *testing.T) {
	req, err := http.NewRequest("GET", "/candidates/1/availability/2", nil)

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

func TestShowCandidateAvailabilityHandlerWhenTimeslotIsForDifferentCandidate(t *testing.T) {
	req, err := http.NewRequest("GET", "/candidates/2/availability/1", nil)

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
	timeSlot := models.TimeSlot{
		Date:      "2020-02-25",
		StartTime: 13,
		EndTime:   16,
	}

	data := []byte(
		fmt.Sprintf(
			`{"date":"%v","start_time":%v,"end_time":%v}`,
			timeSlot.Date,
			timeSlot.StartTime,
			timeSlot.EndTime,
		),
	)

	req, err := http.NewRequest("POST", "/candidates/2/availability", bytes.NewBuffer(data))

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := fmt.Sprintf(
		`{"id":3,"date":"%v","start_time":%v,"end_time":%v}
`,
		timeSlot.Date,
		timeSlot.StartTime,
		timeSlot.EndTime,
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

func TestNewCandidateAvailabilityHandlerMissingFields(t *testing.T) {
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

	req, err := http.NewRequest("POST", "/candidates/2/availability", bytes.NewBuffer(data))

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
		"Missing field \"end_time\" in time_slot",
		resp.Body.String(),
		"Missing field error expected",
	)
}

func TestEditCandidatesAvailabilityHandler(t *testing.T) {
	conf := config.LoadConfig()

	originalTimeSlot := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)[0]

	timeSlot := models.TimeSlot{
		Id:      1,
		EndTime: 14,
	}

	data := []byte(
		fmt.Sprintf(
			`{"start_time":%v,"end_time":%v}`,
			originalTimeSlot["start_time"].(int),
			timeSlot.EndTime,
		),
	)

	req, err := http.NewRequest(
		"PATCH",
		"/candidates/1/availability/1",
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
		`{"id":1,"date":"%v","start_time":%v,"end_time":%v}
`,
		originalTimeSlot["date"].(string),
		originalTimeSlot["start_time"].(int),
		timeSlot.EndTime,
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

func TestDeleteCandidatesAvailabilityHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/candidates/1/availability/1", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	expectedBody := fmt.Sprintf("TimeSlot #%v deleted", 1)

	assert.Equal(t, 200, resp.Code, "200 response expected")

	assert.Equal(
		t,
		expectedBody,
		resp.Body.String(),
		"Deletion message expected",
	)
}

func TestDeleteCandidatesAvailabilityHandlerWhenCandidateDoesNotExist(t *testing.T) {
	conf := config.LoadConfig()
	req, err := http.NewRequest("DELETE", "/candidates/1000/availability/1", nil)

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

func TestDeleteCandidatesAvailabilityHandlerWhenTimeSlotDoesNotExist(t *testing.T) {
	conf := config.LoadConfig()
	req, err := http.NewRequest("DELETE", "/candidates/1/availability/1000", nil)

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
