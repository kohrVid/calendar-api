package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	assert.Equal(t, 200, resp.Code, "200 response expected")
	assert.Equal(t, "application/json; charset=UTF-8", resp.Header().Get("Content-Type"), "JSON response expected")
	assert.Equal(t, "[]\n", resp.Body.String(), "Empty array expected")
}
