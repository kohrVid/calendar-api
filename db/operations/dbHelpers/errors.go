package dbHelpers

import (
	"fmt"
	"strings"

	pluralise "github.com/gertd/go-pluralize"
)

func PgErrorHandler(err error, resource string) string {
	errCode := strings.Split(err.Error(), " ")[1]
	pluralise := pluralise.NewClient()

	switch errCode {
	case "#23505":
		msg := fmt.Sprintf(
			"%v already exists",
			strings.Title(pluralise.Singular(resource)),
		)

		return msg
	default:
		return ""
	}
}
