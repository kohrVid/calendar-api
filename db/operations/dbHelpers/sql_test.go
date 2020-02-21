package dbHelpers

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type Cat struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Colour string `json:"colour"`
}

var cat Cat = Cat{
	Id:     1,
	Name:   "QT",
	Age:    2,
	Colour: "Tabby and white",
}

func TestInsertConfSql(t *testing.T) {
	var conf map[string]interface{}

	dat := fmt.Sprintf(`
database_name: calendar_api_test
data:
  cats:
    - name: %v
      age: %v
      colour: %v`,
		cat.Name,
		cat.Age,
		cat.Colour,
	)

	if err := yaml.Unmarshal([]byte(dat), &conf); err != nil {
		log.Fatal(err)
	}

	assert.Regexp(
		t,
		`INSERT INTO cats \(.*\) VALUES\(.*\);`,
		InsertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
		"formatted",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"name",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"age",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"colour",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		cat.Name,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		strconv.Itoa(cat.Age),
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		cat.Colour,
		"Should return an SQL query with the correct values based on the configuration given",
	)
}

func TestInsertConfSqlWithMultipleRows(t *testing.T) {
	var conf map[string]interface{}

	dat := fmt.Sprintf(`
database_name: calendar_api_test
data:
  cats:
    - name: %v
      age: %v
      colour: %v
    - name: Q-ee
      age: 3
      colour: Grey`,
		cat.Name,
		cat.Age,
		cat.Colour,
	)

	if err := yaml.Unmarshal([]byte(dat), &conf); err != nil {
		log.Fatal(err)
	}

	assert.Regexp(
		t,
		`INSERT INTO cats \(.*\) VALUES\(.*\);`,
		InsertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
		"formatted",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"name",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"age",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"colour",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		cat.Name,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		strconv.Itoa(cat.Age),
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		cat.Colour,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"Q-ee",
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"3",
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"Grey",
		"Should return an SQL query with the correct values based on the configuration given",
	)
}

func TestInsertConfSqlWithMultipleTables(t *testing.T) {
	conf := make(map[string]interface{})

	dat := fmt.Sprintf(`
database_name: calendar_api_test
data:
  cats:
    - name: %v
      age: %v
      colour: %v
  birds:
    - name: Maggie
      colour: Black and white`,
		cat.Name,
		cat.Age,
		cat.Colour,
	)

	if err := yaml.Unmarshal([]byte(dat), &conf); err != nil {
		log.Fatal(err)
	}

	assert.Regexp(
		t,
		`INSERT INTO cats \(.*\) VALUES\(.*\);`,
		InsertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
		"formatted",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"name",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"age",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"colour",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		cat.Name,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		strconv.Itoa(cat.Age),
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		cat.Colour,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Regexp(
		t,
		`INSERT INTO birds \(.*\) VALUES\(.*\);`,
		InsertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
		"formatted",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"Maggie",
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		InsertConfSql(conf),
		"Black and white",
		"Should return an SQL query with the correct values based on the configuration given",
	)
}

func TestInsertConfSqlWithNoTables(t *testing.T) {
	conf := make(map[string]interface{})

	dat := `
database_name: calendar_api_test
data:`

	if err := yaml.Unmarshal([]byte(dat), &conf); err != nil {
		log.Fatal(err)
	}

	assert.Equal(
		t,
		"",
		InsertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
	)
}

func TestSetSqlColumns(t *testing.T) {
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
