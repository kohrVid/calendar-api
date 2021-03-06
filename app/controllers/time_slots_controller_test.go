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

func TestTimeSlotsIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/time_slots", nil)

	if err != nil {
		t.Errorf(
			"Test failed.\nGot:\n\t%v",
			err.Error(),
		)
	}

	resp := httptest.NewRecorder()
	MockRouter().ServeHTTP(resp, req)

	conf := config.LoadConfig()

	timeSlots := config.ToMapList(
		conf["data"].(map[string]interface{})["time_slots"],
	)

	timeSlot1 := timeSlots[0]
	timeSlot2 := timeSlots[1]

	expectedBody := fmt.Sprintf(
		`[{"id":1,"date":"%v","start_time":%v,"end_time":%v},{"id":2,"date":"%v","start_time":%v,"end_time":%v}]
`,
		timeSlot1["date"].(string),
		timeSlot1["start_time"].(int),
		timeSlot1["end_time"].(int),
		timeSlot2["date"].(string),
		timeSlot2["start_time"].(int),
		timeSlot2["end_time"].(int),
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

func TestTimeSlotsIndexHandlerEmpty(t *testing.T) {
	conf := config.LoadConfig()
	dbHelpers.Clean(conf)

	req, err := http.NewRequest("GET", "/time_slots", nil)

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
		"List of time slots expected",
	)

	dbHelpers.Seed(conf)
}

func TestShowTimeSlotsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/time_slots/1", nil)

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
		"JSON of time slot expected",
	)
}

func TestShowTimeSlotsHandlerWhenTimeSlotDoesNotExist(t *testing.T) {
	req, err := http.NewRequest("GET", "/time_slots/1000", nil)

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

func TestNewTimeSlotsHandler(t *testing.T) {
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

	req, err := http.NewRequest("POST", "/time_slots", bytes.NewBuffer(data))

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

func TestNewTimeSlotsHandlerMissingFields(t *testing.T) {
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

	req, err := http.NewRequest("POST", "/time_slots", bytes.NewBuffer(data))

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

func TestEditTimeSlotsHandler(t *testing.T) {
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

	req, err := http.NewRequest("PATCH", "/time_slots/1", bytes.NewBuffer(data))

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

func TestDeleteTimeSlotsHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/time_slots/1", nil)

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

func TestDeleteTimeSlotsHandlerWhenTimeSlotDoesNotExist(t *testing.T) {
	conf := config.LoadConfig()
	req, err := http.NewRequest("DELETE", "/time_slots/1000", nil)

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
