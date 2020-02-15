package dbHelpers

import (
	"fmt"
	"log"

	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
)

func Clean(conf map[string]interface{}) {
	databaseUser := conf["database_user"].(string)

	truncateTables := fmt.Sprintf(`
CREATE OR REPLACE FUNCTION truncate_tables(username IN VARCHAR) RETURNS void AS $$
DECLARE
    statements CURSOR FOR
        SELECT tablename FROM pg_tables
        WHERE tableowner = username AND schemaname = 'public';
BEGIN
    FOR stmt IN statements LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' RESTART IDENTITY CASCADE;';
    END LOOP;
END;
$$ LANGUAGE plpgsql;
		`)

	cleanDB := fmt.Sprintf(
		"%v SELECT truncate_tables('%v');",
		truncateTables,
		databaseUser,
	)

	db := db.DBConnect(conf)
	defer db.Close()

	_, err := db.Exec(cleanDB)
	if err != nil {
		log.Fatal(err)
	}
}

func Seed(conf map[string]interface{}) {
	users := config.ToMapList(conf["users"])
	if len(users) < 1 {
		log.Fatal("No user to seed")
	}

	var seedDB string

	for _, user := range users {
		s := fmt.Sprintf(`
			  INSERT INTO candidates (first_name, last_name, email)
			    VALUES('%v', '%v', '%v');
			`,
			user["first_name"],
			user["last_name"],
			user["email"],
		)
		seedDB += s
	}

	db := db.DBConnect(conf)
	defer db.Close()

	_, err := db.Exec(seedDB)
	if err != nil {
		log.Fatal(err)
	}
}
