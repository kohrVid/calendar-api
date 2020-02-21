package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMapList(t *testing.T) {
	conf := LoadConfig()

	users := conf["data"].(map[string]interface{})["candidates"]

	user1 := map[string]interface{}{
		"first_name": "Alex Courtney",
		"last_name":  "Dusk",
		"email":      "alex.c.dusk@example.com",
	}

	user2 := map[string]interface{}{
		"first_name": "Aisha",
		"last_name":  "Prince",
		"email":      "aisha.prince@example.com",
	}

	expected := []map[string]interface{}{user1, user2}

	assert.Equal(t, expected, ToMapList(users), "expected to return a list of map")
}

func TestToMapListWithEmptyList(t *testing.T) {
	conf := LoadConfig()
	users := conf["no_users"]
	expected := []map[string]interface{}{}
	assert.Equal(
		t,
		expected,
		ToMapList(users),
		"expected to return an empty list",
	)
}
