package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kohrVid/calendar-api/config"
	"github.com/stretchr/testify/assert"
)

func TestCandidatesHandler(t *testing.T) {
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
