package dbHelpers

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	pluralise "github.com/gertd/go-pluralize"
	"github.com/stretchr/testify/assert"
)

func TestPgErrorHandlerNil(t *testing.T) {
	resource := "candidates"
	err := errors.New("Non postgres error message")

	assert.Equal(
		t,
		"",
		PgErrorHandler(err, resource),
		"Should return nil for non postgres error messages",
	)
}

func TestPgErrorHandlerDuplicate(t *testing.T) {

	pluralise := pluralise.NewClient()
	resource := "candidates"

	err := errors.New(
		"ERROR #23505 duplicate key value violates unique constraint \"candidates_unique_idx\"",
	)

	assert.Equal(
		t,
		fmt.Sprintf(
			"%v already exists",
			strings.Title(pluralise.Singular(resource)),
		),
		PgErrorHandler(err, resource),
		"Should return duplicate row error messages for resource",
	)
}

func TestPgErrorHandlerMissingFields(t *testing.T) {
	pluralise := pluralise.NewClient()
	column := "first_name"
	resource := "candidates"

	err := errors.New(
		fmt.Sprintf(
			"ERROR #23502 null value in column \"%v\" violates not-null constraint",
			column,
		),
	)

	assert.Equal(
		t,
		fmt.Sprintf(
			"Missing field \"%v\" in %v",
			column,
			pluralise.Singular(resource),
		),
		PgErrorHandler(err, resource),
		"Should return missing fields error messages for resource",
	)
}
