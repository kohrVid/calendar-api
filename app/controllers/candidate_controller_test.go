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
	users := config.ToMapList(conf["users"])
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
	user := config.ToMapList(conf["users"])[0]

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
