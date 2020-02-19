package dbHelpers

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

type Cat struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Colour string `json:"colour"`
}

func TestSetSqlColumns(t *testing.T) {
	cat := Cat{
		Id:     1,
		Name:   "QT",
		Age:    2,
		Colour: "Tabby and white",
	}

	params := Cat{
		Name: "Q-ee",
	}

	c := structs.New(&cat)
	p := structs.New(&params)

	assert.Equal(
		t,
		"SET name = 'Q-ee'",
		SetSqlColumns(c, p),
		"Should return the correct SQL query if a single column is changed",
	)
}

func TestSetSqlColumnsWithMultipleColumns(t *testing.T) {
	cat := Cat{
		Id:     1,
		Name:   "QT",
		Age:    2,
		Colour: "Tabby and white",
	}

	params := Cat{
		Name:   "Q-ee",
		Colour: "blue",
	}

	c := structs.New(&cat)
	p := structs.New(&params)

	assert.Equal(
		t,
		"SET name = 'Q-ee', colour = 'blue'",
		SetSqlColumns(c, p),
		"Should return the correct SQL query if multiple columns are changed",
	)
}

func TestSetSqlColumnsWithNonStrings(t *testing.T) {
	cat := Cat{
		Id:     1,
		Name:   "QT",
		Age:    2,
		Colour: "Tabby and white",
	}

	params := Cat{
		Id: 3,
	}

	c := structs.New(&cat)
	p := structs.New(&params)

	assert.Equal(
		t,
		"SET id = 3",
		SetSqlColumns(c, p),
		"Should return the correct SQL query if multiple columns are changed",
	)
}
