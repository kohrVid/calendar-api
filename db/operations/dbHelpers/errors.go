package dbHelpers

import (
	"fmt"
	"strings"

	pluralise "github.com/gertd/go-pluralize"
)

func PgErrorHandler(err error, resource string) string {
	e := strings.Split(err.Error(), " ")
	errCode := e[1]
	pluralise := pluralise.NewClient()
	msg := ""

	switch errCode {
	case "#23502":
		msg = fmt.Sprintf(
			"Missing field %v in %v",
			e[6],
			pluralise.Singular(resource),
		)

	case "#23505":
		msg = fmt.Sprintf(
			"%v already exists",
			strings.Title(pluralise.Singular(resource)),
		)

	default:

	}

	return msg
}
