package db

import (
	"crypto/tls"
	"os"

	"github.com/go-pg/pg"
)

func DBConnect(config map[string]interface{}) *pg.DB {
	db := pg.Connect(&pg.Options{
		User:      config["database_user"].(string),
		Database:  config["database_name"].(string),
		TLSConfig: sslMode(),
	})

	return db
}

func sslMode() *tls.Config {
	switch os.Getenv("SSL_MODE") {
	case "verify-ca", "verify-full":
		return &tls.Config{}
	case "allow", "prefer", "require":
		return &tls.Config{InsecureSkipVerify: true}
	default:
		return nil
	}
}
