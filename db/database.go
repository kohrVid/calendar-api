package db

import (
	"crypto/tls"
	"os"

	"github.com/go-pg/pg"
)

func DBConnect() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:      "calendar_api",
		Database:  "calendar_api",
		TLSConfig: sslMode(),
	})

	return db
}

func sslMode() *tls.Config {
	switch os.Getenv("SSL_MODE") {
	case "verify-ca", "verify-full":
		return &tls.Config{}
	case "allow", "prefer", "require":
		return &tls.Config{InsecureSkipVerify: true} //nolint
	default:
		return nil
	}
}
